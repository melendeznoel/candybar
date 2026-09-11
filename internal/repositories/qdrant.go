package repositories

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

// defaultQdrantPort is the gRPC port Qdrant listens on (the REST API is 6333).
const defaultQdrantPort = 6334

// defaultQdrantTimeout is applied to any request whose context carries no
// deadline of its own.
const defaultQdrantTimeout = 30 * time.Second

// Qdrant is a reusable wrapper around the official Qdrant Go client. It keeps
// the verbose protobuf request building in one place and exposes the handful of
// operations feature packages actually need, in plain Go types.
//
// A Qdrant value is safe for concurrent use and is meant to be created once and
// shared, the same way a *sql.DB is. Call Close when the process shuts down.
type Qdrant struct {
	client  *qdrant.Client
	timeout time.Duration
}

// QdrantPoint is a single vector plus the payload stored alongside it. ID must
// be an unsigned integer or a UUID string - the only two forms Qdrant accepts.
type QdrantPoint struct {
	ID      interface{}
	Vector  []float32
	Payload map[string]interface{}
}

// QdrantSearchResult is one point returned by Search, Retrieve or Scroll.
// Score is zero for results that did not come from a similarity query.
type QdrantSearchResult struct {
	ID      interface{}
	Score   float32
	Version uint64
	Vector  []float32
	Payload map[string]interface{}
}

// QdrantSearchRequest describes a nearest-neighbour query. Filter and Params are
// the client library types, so any condition the server understands can be built
// with the qdrant.NewMatch* helpers and passed straight through.
type QdrantSearchRequest struct {
	Vector         []float32
	Limit          uint64
	Offset         uint64
	Filter         *qdrant.Filter
	Params         *qdrant.SearchParams
	ScoreThreshold *float32
	WithPayload    bool
	WithVector     bool
}

// NewQdrant dials a Qdrant instance over gRPC. apiKey may be empty for an
// unsecured local instance; useTLS should be true for Qdrant Cloud.
func NewQdrant(host string, port int, apiKey string, useTLS bool) (*Qdrant, error) {
	if host == "" {
		host = "localhost"
	}

	if port == 0 {
		port = defaultQdrantPort
	}

	return NewQdrantWithConfig(&qdrant.Config{
		Host:   host,
		Port:   port,
		APIKey: apiKey,
		UseTLS: useTLS,
	})
}

// NewQdrantWithConfig builds a client from a full qdrant.Config, for callers
// that need connection pooling, keepalives, retries or custom TLS.
func NewQdrantWithConfig(config *qdrant.Config) (*Qdrant, error) {
	client, err := qdrant.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("could not connect to qdrant: %s", err)
	}

	return &Qdrant{client: client, timeout: defaultQdrantTimeout}, nil
}

// QdrantFromEnv builds a client from the environment.
//
//	export QDRANT_URL="https://[CLUSTER_ID].[REGION].cloud.qdrant.io:6334"
//	# or, instead of QDRANT_URL:
//	export QDRANT_HOST="localhost"
//	export QDRANT_PORT="6334"
//	export QDRANT_API_KEY="[API_KEY]"
//
// A QDRANT_URL with an https scheme turns TLS on; otherwise set QDRANT_USE_TLS.
func QdrantFromEnv() (*Qdrant, error) {
	host := os.Getenv("QDRANT_HOST")
	port := defaultQdrantPort
	useTLS := strings.EqualFold(os.Getenv("QDRANT_USE_TLS"), "true")

	if rawURL := os.Getenv("QDRANT_URL"); rawURL != "" {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return nil, fmt.Errorf("could not parse QDRANT_URL: %s", err)
		}

		host = parsed.Hostname()
		useTLS = useTLS || parsed.Scheme == "https"

		if parsed.Port() != "" {
			if port, err = strconv.Atoi(parsed.Port()); err != nil {
				return nil, fmt.Errorf("could not parse port in QDRANT_URL: %s", err)
			}
		}
	} else if envPort := os.Getenv("QDRANT_PORT"); envPort != "" {
		parsedPort, err := strconv.Atoi(envPort)
		if err != nil {
			return nil, fmt.Errorf("could not parse QDRANT_PORT: %s", err)
		}

		port = parsedPort
	}

	if host == "" {
		return nil, fmt.Errorf("neither QDRANT_URL nor QDRANT_HOST is set")
	}

	return NewQdrant(host, port, os.Getenv("QDRANT_API_KEY"), useTLS)
}

// WithTimeout overrides the deadline applied to requests whose context has none,
// and returns the receiver for chaining. A non-positive duration disables it.
func (q *Qdrant) WithTimeout(timeout time.Duration) *Qdrant {
	q.timeout = timeout

	return q
}

// Client exposes the underlying qdrant.Client for the operations this wrapper
// does not cover (snapshots, aliases, batch updates, grouped queries).
func (q *Qdrant) Client() *qdrant.Client {
	return q.client
}

// Close releases the gRPC connection.
func (q *Qdrant) Close() error {
	return q.client.Close()
}

// Healthy reports whether the Qdrant instance answers a health check.
func (q *Qdrant) Healthy(ctx context.Context) error {
	ctx, cancel := q.context(ctx)
	defer cancel()

	if _, err := q.client.HealthCheck(ctx); err != nil {
		return fmt.Errorf("qdrant health check failed: %s", err)
	}

	return nil
}

// CreateCollection creates a collection holding vectors of the given size and
// distance metric (for example qdrant.Distance_Cosine). It is a no-op when the
// collection already exists.
func (q *Qdrant) CreateCollection(ctx context.Context, collection string, vectorSize uint64, distance qdrant.Distance) error {
	if vectorSize == 0 {
		return fmt.Errorf("vector size must be greater than zero")
	}

	if distance == qdrant.Distance_UnknownDistance {
		distance = qdrant.Distance_Cosine
	}

	exists, err := q.CollectionExists(ctx, collection)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	ctx, cancel := q.context(ctx)
	defer cancel()

	request := &qdrant.CreateCollection{
		CollectionName: collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vectorSize,
			Distance: distance,
		}),
	}

	if err := q.client.CreateCollection(ctx, request); err != nil {
		return fmt.Errorf("could not create collection %s: %s", collection, err)
	}

	return nil
}

// CollectionExists reports whether a collection is present on the server.
func (q *Qdrant) CollectionExists(ctx context.Context, collection string) (bool, error) {
	ctx, cancel := q.context(ctx)
	defer cancel()

	exists, err := q.client.CollectionExists(ctx, collection)
	if err != nil {
		return false, fmt.Errorf("could not check collection %s: %s", collection, err)
	}

	return exists, nil
}

// DeleteCollection drops a collection and every point in it.
func (q *Qdrant) DeleteCollection(ctx context.Context, collection string) error {
	ctx, cancel := q.context(ctx)
	defer cancel()

	if err := q.client.DeleteCollection(ctx, collection); err != nil {
		return fmt.Errorf("could not delete collection %s: %s", collection, err)
	}

	return nil
}

// ListCollections returns the names of every collection on the server.
func (q *Qdrant) ListCollections(ctx context.Context) ([]string, error) {
	ctx, cancel := q.context(ctx)
	defer cancel()

	collections, err := q.client.ListCollections(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list collections: %s", err)
	}

	return collections, nil
}

// Upsert inserts or replaces points, waiting for the write to be applied so a
// subsequent Search sees them.
func (q *Qdrant) Upsert(ctx context.Context, collection string, points []QdrantPoint) error {
	if len(points) == 0 {
		return nil
	}

	structs := make([]*qdrant.PointStruct, 0, len(points))

	for _, point := range points {
		id, err := QdrantID(point.ID)
		if err != nil {
			return err
		}

		if len(point.Vector) == 0 {
			return fmt.Errorf("point %v has an empty vector", point.ID)
		}

		payload, err := qdrant.TryValueMap(point.Payload)
		if err != nil {
			return fmt.Errorf("could not encode payload for point %v: %s", point.ID, err)
		}

		structs = append(structs, &qdrant.PointStruct{
			Id:      id,
			Vectors: qdrant.NewVectorsDense(point.Vector),
			Payload: payload,
		})
	}

	ctx, cancel := q.context(ctx)
	defer cancel()

	request := &qdrant.UpsertPoints{
		CollectionName: collection,
		Points:         structs,
		Wait:           qdrant.PtrOf(true),
	}

	if _, err := q.client.Upsert(ctx, request); err != nil {
		return fmt.Errorf("could not upsert into %s: %s", collection, err)
	}

	return nil
}

// Search returns the points nearest to request.Vector, ordered by descending
// score.
func (q *Qdrant) Search(ctx context.Context, collection string, request QdrantSearchRequest) ([]QdrantSearchResult, error) {
	if len(request.Vector) == 0 {
		return nil, fmt.Errorf("search vector is empty")
	}

	if request.Limit == 0 {
		request.Limit = 10
	}

	ctx, cancel := q.context(ctx)
	defer cancel()

	query := &qdrant.QueryPoints{
		CollectionName: collection,
		Query:          qdrant.NewQueryDense(request.Vector),
		Limit:          qdrant.PtrOf(request.Limit),
		Filter:         request.Filter,
		Params:         request.Params,
		ScoreThreshold: request.ScoreThreshold,
		WithPayload:    qdrant.NewWithPayload(request.WithPayload),
		WithVectors:    qdrant.NewWithVectors(request.WithVector),
	}

	if request.Offset > 0 {
		query.Offset = qdrant.PtrOf(request.Offset)
	}

	points, err := q.client.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not search %s: %s", collection, err)
	}

	results := make([]QdrantSearchResult, 0, len(points))

	for _, point := range points {
		results = append(results, QdrantSearchResult{
			ID:      qdrantIDValue(point.GetId()),
			Score:   point.GetScore(),
			Version: point.GetVersion(),
			Vector:  qdrantVectorValues(point.GetVectors()),
			Payload: qdrantPayloadMap(point.GetPayload()),
		})
	}

	return results, nil
}

// Retrieve fetches points by ID, including their payload and vector.
func (q *Qdrant) Retrieve(ctx context.Context, collection string, ids []interface{}) ([]QdrantSearchResult, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	pointIDs, err := qdrantIDs(ids)
	if err != nil {
		return nil, err
	}

	ctx, cancel := q.context(ctx)
	defer cancel()

	request := &qdrant.GetPoints{
		CollectionName: collection,
		Ids:            pointIDs,
		WithPayload:    qdrant.NewWithPayload(true),
		WithVectors:    qdrant.NewWithVectors(true),
	}

	points, err := q.client.Get(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve from %s: %s", collection, err)
	}

	return qdrantRetrievedResults(points), nil
}

// Delete removes points by ID.
func (q *Qdrant) Delete(ctx context.Context, collection string, ids []interface{}) error {
	if len(ids) == 0 {
		return nil
	}

	pointIDs, err := qdrantIDs(ids)
	if err != nil {
		return err
	}

	return q.delete(ctx, collection, qdrant.NewPointsSelectorIDs(pointIDs))
}

// DeleteByFilter removes every point matching a filter. A nil filter is rejected
// so a bug cannot silently empty a collection - use DeleteCollection for that.
func (q *Qdrant) DeleteByFilter(ctx context.Context, collection string, filter *qdrant.Filter) error {
	if filter == nil {
		return fmt.Errorf("filter is nil: refusing to delete every point in %s", collection)
	}

	return q.delete(ctx, collection, qdrant.NewPointsSelectorFilter(filter))
}

// Count returns an exact count of the points matching filter. A nil filter counts
// the whole collection.
func (q *Qdrant) Count(ctx context.Context, collection string, filter *qdrant.Filter) (uint64, error) {
	ctx, cancel := q.context(ctx)
	defer cancel()

	request := &qdrant.CountPoints{
		CollectionName: collection,
		Filter:         filter,
		Exact:          qdrant.PtrOf(true),
	}

	count, err := q.client.Count(ctx, request)
	if err != nil {
		return 0, fmt.Errorf("could not count points in %s: %s", collection, err)
	}

	return count, nil
}

// Scroll pages through a collection, returning up to limit points along with the
// offset to pass back in for the next page. A nil next offset means the last page
// was reached.
func (q *Qdrant) Scroll(ctx context.Context, collection string, limit uint32, offset interface{}, filter *qdrant.Filter) ([]QdrantSearchResult, interface{}, error) {
	if limit == 0 {
		limit = 100
	}

	request := &qdrant.ScrollPoints{
		CollectionName: collection,
		Filter:         filter,
		Limit:          qdrant.PtrOf(limit),
		WithPayload:    qdrant.NewWithPayload(true),
		WithVectors:    qdrant.NewWithVectors(false),
	}

	if offset != nil {
		startAt, err := QdrantID(offset)
		if err != nil {
			return nil, nil, err
		}

		request.Offset = startAt
	}

	ctx, cancel := q.context(ctx)
	defer cancel()

	points, next, err := q.client.ScrollAndOffset(ctx, request)
	if err != nil {
		return nil, nil, fmt.Errorf("could not scroll %s: %s", collection, err)
	}

	if next == nil {
		return qdrantRetrievedResults(points), nil, nil
	}

	return qdrantRetrievedResults(points), qdrantIDValue(next), nil
}

// delete runs a delete against the given selector, waiting for it to be applied.
func (q *Qdrant) delete(ctx context.Context, collection string, selector *qdrant.PointsSelector) error {
	ctx, cancel := q.context(ctx)
	defer cancel()

	request := &qdrant.DeletePoints{
		CollectionName: collection,
		Points:         selector,
		Wait:           qdrant.PtrOf(true),
	}

	if _, err := q.client.Delete(ctx, request); err != nil {
		return fmt.Errorf("could not delete from %s: %s", collection, err)
	}

	return nil
}

// context applies the client timeout to a context that does not already carry a
// deadline of its own.
func (q *Qdrant) context(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}

	if q.timeout <= 0 {
		return context.WithCancel(ctx)
	}

	if _, ok := ctx.Deadline(); ok {
		return context.WithCancel(ctx)
	}

	return context.WithTimeout(ctx, q.timeout)
}

// QdrantID converts a Go value to a Qdrant point ID. Qdrant accepts only
// unsigned integers and UUID strings.
func QdrantID(id interface{}) (*qdrant.PointId, error) {
	switch value := id.(type) {
	case *qdrant.PointId:
		return value, nil
	case uint64:
		return qdrant.NewIDNum(value), nil
	case uint:
		return qdrant.NewIDNum(uint64(value)), nil
	case uint32:
		return qdrant.NewIDNum(uint64(value)), nil
	case int:
		if value < 0 {
			return nil, fmt.Errorf("point id %d is negative", value)
		}

		return qdrant.NewIDNum(uint64(value)), nil
	case int64:
		if value < 0 {
			return nil, fmt.Errorf("point id %d is negative", value)
		}

		return qdrant.NewIDNum(uint64(value)), nil
	case string:
		if value == "" {
			return nil, fmt.Errorf("point id is empty")
		}

		return qdrant.NewIDUUID(value), nil
	default:
		return nil, fmt.Errorf("unsupported point id type %T: use an unsigned integer or a UUID string", id)
	}
}

// qdrantIDs converts a slice of Go values to Qdrant point IDs.
func qdrantIDs(ids []interface{}) ([]*qdrant.PointId, error) {
	pointIDs := make([]*qdrant.PointId, 0, len(ids))

	for _, id := range ids {
		pointID, err := QdrantID(id)
		if err != nil {
			return nil, err
		}

		pointIDs = append(pointIDs, pointID)
	}

	return pointIDs, nil
}

// qdrantRetrievedResults maps points returned by Get and Scroll onto the plain
// result type. Those responses carry no similarity score.
func qdrantRetrievedResults(points []*qdrant.RetrievedPoint) []QdrantSearchResult {
	results := make([]QdrantSearchResult, 0, len(points))

	for _, point := range points {
		results = append(results, QdrantSearchResult{
			ID:      qdrantIDValue(point.GetId()),
			Vector:  qdrantVectorValues(point.GetVectors()),
			Payload: qdrantPayloadMap(point.GetPayload()),
		})
	}

	return results
}

// qdrantIDValue unwraps a point ID into a uint64 or a UUID string.
func qdrantIDValue(id *qdrant.PointId) interface{} {
	if id == nil {
		return nil
	}

	if uuid := id.GetUuid(); uuid != "" {
		return uuid
	}

	return id.GetNum()
}

// qdrantVectorValues pulls the dense vector out of a response, returning nil for
// sparse, multi or named vectors, which this wrapper does not model.
func qdrantVectorValues(vectors *qdrant.VectorsOutput) []float32 {
	if vectors == nil {
		return nil
	}

	return vectors.GetVector().GetData()
}

// qdrantPayloadMap converts a Qdrant payload back into a plain Go map.
func qdrantPayloadMap(payload map[string]*qdrant.Value) map[string]interface{} {
	if len(payload) == 0 {
		return nil
	}

	result := make(map[string]interface{}, len(payload))

	for key, value := range payload {
		result[key] = qdrantValue(value)
	}

	return result
}

// qdrantValue converts a single payload value back into a plain Go value.
func qdrantValue(value *qdrant.Value) interface{} {
	if value == nil {
		return nil
	}

	switch value.GetKind().(type) {
	case *qdrant.Value_NullValue:
		return nil
	case *qdrant.Value_BoolValue:
		return value.GetBoolValue()
	case *qdrant.Value_IntegerValue:
		return value.GetIntegerValue()
	case *qdrant.Value_DoubleValue:
		return value.GetDoubleValue()
	case *qdrant.Value_StringValue:
		return value.GetStringValue()
	case *qdrant.Value_StructValue:
		return qdrantPayloadMap(value.GetStructValue().GetFields())
	case *qdrant.Value_ListValue:
		items := value.GetListValue().GetValues()
		list := make([]interface{}, 0, len(items))

		for _, item := range items {
			list = append(list, qdrantValue(item))
		}

		return list
	default:
		return nil
	}
}
