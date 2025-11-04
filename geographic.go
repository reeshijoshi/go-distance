package distance

import "math"

const (
	earthRadiusKm    = 6371.0 // Earth's mean radius in kilometers
	earthRadiusMiles = 3959.0 // Earth's mean radius in miles
	degToRad         = math.Pi / 180.0
)

// Coord represents a geographic coordinate (latitude, longitude).
type Coord struct {
	Lat float64 // Latitude in degrees [-90, 90]
	Lon float64 // Longitude in degrees [-180, 180]
}

// Haversine computes great-circle distance using Haversine formula.
// Returns distance in kilometers by default.
// Time: O(1), Space: O(1)
func Haversine(a, b Coord) float64 {
	return HaversineWithRadius(a, b, earthRadiusKm)
}

// HaversineWithRadius computes Haversine distance with custom Earth radius.
// Time: O(1), Space: O(1)
func HaversineWithRadius(a, b Coord, radius float64) float64 {
	lat1 := a.Lat * degToRad
	lat2 := b.Lat * degToRad
	deltaLat := (b.Lat - a.Lat) * degToRad
	deltaLon := (b.Lon - a.Lon) * degToRad

	// Haversine formula
	sinDLat := math.Sin(deltaLat / 2)
	sinDLon := math.Sin(deltaLon / 2)
	h := sinDLat*sinDLat + math.Cos(lat1)*math.Cos(lat2)*sinDLon*sinDLon

	c := 2 * math.Atan2(math.Sqrt(h), math.Sqrt(1-h))

	return radius * c
}

// HaversineMiles computes Haversine distance in miles.
// Time: O(1), Space: O(1)
func HaversineMiles(a, b Coord) float64 {
	return HaversineWithRadius(a, b, earthRadiusMiles)
}

// GreatCircle computes great-circle distance using spherical law of cosines.
// Simpler but less accurate for small distances than Haversine.
// Returns distance in kilometers.
// Time: O(1), Space: O(1)
func GreatCircle(a, b Coord) float64 {
	return GreatCircleWithRadius(a, b, earthRadiusKm)
}

// GreatCircleWithRadius computes great-circle distance with custom radius.
// Time: O(1), Space: O(1)
func GreatCircleWithRadius(a, b Coord, radius float64) float64 {
	lat1 := a.Lat * degToRad
	lat2 := b.Lat * degToRad
	deltaLon := (b.Lon - a.Lon) * degToRad

	// Spherical law of cosines
	cosAngle := math.Sin(lat1)*math.Sin(lat2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Cos(deltaLon)

	// Clamp to [-1, 1] for numerical stability
	if cosAngle > 1 {
		cosAngle = 1
	} else if cosAngle < -1 {
		cosAngle = -1
	}

	angle := math.Acos(cosAngle)
	return radius * angle
}

// Equirectangular computes approximate distance using equirectangular projection.
// Fast but less accurate approximation for small distances.
// Returns distance in kilometers.
// Time: O(1), Space: O(1)
func Equirectangular(a, b Coord) float64 {
	return EquirectangularWithRadius(a, b, earthRadiusKm)
}

// EquirectangularWithRadius computes equirectangular distance with custom radius.
// Time: O(1), Space: O(1)
func EquirectangularWithRadius(a, b Coord, radius float64) float64 {
	lat1 := a.Lat * degToRad
	lat2 := b.Lat * degToRad
	deltaLon := (b.Lon - a.Lon) * degToRad
	deltaLat := (b.Lat - a.Lat) * degToRad

	x := deltaLon * math.Cos((lat1+lat2)/2)
	y := deltaLat

	return radius * math.Sqrt(x*x+y*y)
}

// Vincenty computes geodesic distance using Vincenty formula.
// More accurate than Haversine for oblate spheroid (WGS-84 ellipsoid).
// Returns distance in meters.
// Time: O(1) with iteration, Space: O(1)
func Vincenty(a, b Coord) (float64, error) {
	const (
		majorAxis     = 6378137.0         // WGS-84 semi-major axis (meters)
		minorAxis     = 6356752.314245    // WGS-84 semi-minor axis (meters)
		flattening    = 1 / 298.257223563 // WGS-84 flattening
		tolerance     = 1e-12
		maxIterations = 200
	)

	lat1 := a.Lat * degToRad
	lat2 := b.Lat * degToRad
	lon1 := a.Lon * degToRad
	lon2 := b.Lon * degToRad

	L := lon2 - lon1

	U1 := math.Atan((1 - flattening) * math.Tan(lat1))
	U2 := math.Atan((1 - flattening) * math.Tan(lat2))

	sinU1, cosU1 := math.Sin(U1), math.Cos(U1)
	sinU2, cosU2 := math.Sin(U2), math.Cos(U2)

	lambda := L
	var lambdaP float64

	var sinSigma, cosSigma, sigma, sinAlpha, cosSqAlpha, cos2SigmaM float64
	converged := false

	for i := 0; i < maxIterations; i++ {
		sinLambda, cosLambda := math.Sin(lambda), math.Cos(lambda)

		sinSigma = math.Sqrt(
			(cosU2*sinLambda)*(cosU2*sinLambda) +
				(cosU1*sinU2-sinU1*cosU2*cosLambda)*(cosU1*sinU2-sinU1*cosU2*cosLambda),
		)

		if sinSigma == 0 {
			return 0, nil // Coincident points
		}

		cosSigma = sinU1*sinU2 + cosU1*cosU2*cosLambda
		sigma = math.Atan2(sinSigma, cosSigma)
		sinAlpha = cosU1 * cosU2 * sinLambda / sinSigma
		cosSqAlpha = 1 - sinAlpha*sinAlpha
		cos2SigmaM = cosSigma - 2*sinU1*sinU2/cosSqAlpha

		if math.IsNaN(cos2SigmaM) {
			cos2SigmaM = 0 // Equatorial line
		}

		C := flattening / 16 * cosSqAlpha * (4 + flattening*(4-3*cosSqAlpha))

		lambdaP = lambda
		lambda = L + (1-C)*flattening*sinAlpha*
			(sigma+C*sinSigma*(cos2SigmaM+C*cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)))

		if math.Abs(lambda-lambdaP) < tolerance {
			converged = true
			break
		}
	}

	// Check if algorithm converged
	if !converged {
		// For antipodal points or nearly antipodal points, formula may not converge
		// Fall back to Haversine as approximation
		return HaversineWithRadius(a, b, majorAxis/1000.0) * 1000.0, nil
	}

	uSq := cosSqAlpha * (majorAxis*majorAxis - minorAxis*minorAxis) / (minorAxis * minorAxis)
	A := 1 + uSq/16384*(4096+uSq*(-768+uSq*(320-175*uSq)))
	B := uSq / 1024 * (256 + uSq*(-128+uSq*(74-47*uSq)))

	deltaSigma := B * sinSigma * (cos2SigmaM + B/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-
		B/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM)))

	s := minorAxis * A * (sigma - deltaSigma)

	return s, nil
}

// VincentyKm computes Vincenty distance in kilometers.
// Time: O(1) with iteration, Space: O(1)
func VincentyKm(a, b Coord) (float64, error) {
	meters, err := Vincenty(a, b)
	if err != nil {
		return 0, err
	}
	return meters / 1000.0, nil
}
