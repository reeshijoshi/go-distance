package distance

import (
	"math"
	"testing"
)

func TestHaversine(t *testing.T) {
	tests := []struct {
		name      string
		a, b      Coord
		expected  float64
		tolerance float64
	}{
		{
			name:      "same location",
			a:         Coord{Lat: 40.7128, Lon: -74.0060},
			b:         Coord{Lat: 40.7128, Lon: -74.0060},
			expected:  0,
			tolerance: 0.1,
		},
		{
			name:      "NYC to London",
			a:         Coord{Lat: 40.7128, Lon: -74.0060},
			b:         Coord{Lat: 51.5074, Lon: -0.1278},
			expected:  5570,
			tolerance: 10,
		},
		{
			name:      "San Francisco to Tokyo",
			a:         Coord{Lat: 37.7749, Lon: -122.4194},
			b:         Coord{Lat: 35.6762, Lon: 139.6503},
			expected:  8280,
			tolerance: 50,
		},
		{
			name:      "Sydney to Cape Town",
			a:         Coord{Lat: -33.8688, Lon: 151.2093},
			b:         Coord{Lat: -33.9249, Lon: 18.4241},
			expected:  11000,
			tolerance: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Haversine(tt.a, tt.b)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("expected %v km (±%v), got %v km", tt.expected, tt.tolerance, result)
			}
		})
	}
}

func TestHaversineMiles(t *testing.T) {
	nyc := Coord{Lat: 40.7128, Lon: -74.0060}
	london := Coord{Lat: 51.5074, Lon: -0.1278}

	result := HaversineMiles(nyc, london)

	// NYC to London is approximately 3461 miles
	expected := 3461.0
	tolerance := 10.0

	if math.Abs(result-expected) > tolerance {
		t.Errorf("expected %v miles (±%v), got %v miles", expected, tolerance, result)
	}
}

func TestHaversineWithRadius(t *testing.T) {
	a := Coord{Lat: 0, Lon: 0}
	b := Coord{Lat: 0, Lon: 90}

	// Quarter of circle on equator with custom radius
	radius := 1000.0
	result := HaversineWithRadius(a, b, radius)

	// Expected: quarter circumference = 2*pi*r/4 = pi*r/2
	expected := math.Pi * radius / 2
	tolerance := 1.0

	if math.Abs(result-expected) > tolerance {
		t.Errorf("expected %v (±%v), got %v", expected, tolerance, result)
	}
}

func TestGreatCircle(t *testing.T) {
	tests := []struct {
		name      string
		a, b      Coord
		expected  float64
		tolerance float64
	}{
		{
			name:      "same location",
			a:         Coord{Lat: 0, Lon: 0},
			b:         Coord{Lat: 0, Lon: 0},
			expected:  0,
			tolerance: 0.1,
		},
		{
			name:      "equator quarter",
			a:         Coord{Lat: 0, Lon: 0},
			b:         Coord{Lat: 0, Lon: 90},
			expected:  10018,
			tolerance: 50,
		},
		{
			name:      "pole to pole",
			a:         Coord{Lat: 90, Lon: 0},
			b:         Coord{Lat: -90, Lon: 0},
			expected:  20015,
			tolerance: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GreatCircle(tt.a, tt.b)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("expected %v km (±%v), got %v km", tt.expected, tt.tolerance, result)
			}
		})
	}
}

func TestEquirectangular(t *testing.T) {
	tests := []struct {
		name      string
		a, b      Coord
		expected  float64
		tolerance float64
	}{
		{
			name:      "same location",
			a:         Coord{Lat: 40, Lon: -74},
			b:         Coord{Lat: 40, Lon: -74},
			expected:  0,
			tolerance: 0.1,
		},
		{
			name:      "small distance",
			a:         Coord{Lat: 40.0, Lon: -74.0},
			b:         Coord{Lat: 40.1, Lon: -74.1},
			expected:  13,
			tolerance: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Equirectangular(tt.a, tt.b)
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("expected %v km (±%v), got %v km", tt.expected, tt.tolerance, result)
			}
		})
	}
}

func TestVincenty(t *testing.T) {
	tests := []struct {
		name      string
		a, b      Coord
		expected  float64
		tolerance float64
	}{
		{
			name:      "same location",
			a:         Coord{Lat: 40.7128, Lon: -74.0060},
			b:         Coord{Lat: 40.7128, Lon: -74.0060},
			expected:  0,
			tolerance: 1,
		},
		{
			name:      "NYC to London (meters)",
			a:         Coord{Lat: 40.7128, Lon: -74.0060},
			b:         Coord{Lat: 51.5074, Lon: -0.1278},
			expected:  5570000,
			tolerance: 20000,
		},
		{
			name:      "short distance",
			a:         Coord{Lat: 40.0, Lon: -74.0},
			b:         Coord{Lat: 40.1, Lon: -74.0},
			expected:  11119,
			tolerance: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Vincenty(tt.a, tt.b)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > tt.tolerance {
				t.Errorf("expected %v m (±%v), got %v m", tt.expected, tt.tolerance, result)
			}
		})
	}
}

func TestVincentyKm(t *testing.T) {
	nyc := Coord{Lat: 40.7128, Lon: -74.0060}
	london := Coord{Lat: 51.5074, Lon: -0.1278}

	result, err := VincentyKm(nyc, london)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := 5570.0
	tolerance := 20.0

	if math.Abs(result-expected) > tolerance {
		t.Errorf("expected %v km (±%v), got %v km", expected, tolerance, result)
	}
}

func TestGeographicEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		a, b Coord
	}{
		{
			name: "crosses international date line",
			a:    Coord{Lat: 0, Lon: -179},
			b:    Coord{Lat: 0, Lon: 179},
		},
		{
			name: "near north pole",
			a:    Coord{Lat: 89, Lon: 0},
			b:    Coord{Lat: 89, Lon: 180},
		},
		{
			name: "near south pole",
			a:    Coord{Lat: -89, Lon: 0},
			b:    Coord{Lat: -89, Lon: 180},
		},
		{
			name: "crosses equator",
			a:    Coord{Lat: -10, Lon: 0},
			b:    Coord{Lat: 10, Lon: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that all methods complete without error
			h := Haversine(tt.a, tt.b)
			g := GreatCircle(tt.a, tt.b)
			e := Equirectangular(tt.a, tt.b)
			v, err := Vincenty(tt.a, tt.b)

			if err != nil {
				t.Errorf("Vincenty error: %v", err)
			}

			// All distances should be positive
			if h < 0 || g < 0 || e < 0 || v < 0 {
				t.Errorf("negative distance detected: h=%v, g=%v, e=%v, v=%v", h, g, e, v)
			}

			// Haversine and Vincenty should be reasonably close (within 1%)
			if v > 0 {
				vKm := v / 1000
				diff := math.Abs(h-vKm) / vKm
				if diff > 0.01 {
					t.Logf("Haversine vs Vincenty difference: %.2f%% (h=%v, v=%v km)", diff*100, h, vKm)
				}
			}
		})
	}
}

func TestGeographicConsistency(t *testing.T) {
	// Test that distance(a,b) == distance(b,a)
	coords := []Coord{
		{Lat: 40.7128, Lon: -74.0060},  // NYC
		{Lat: 51.5074, Lon: -0.1278},   // London
		{Lat: 35.6762, Lon: 139.6503},  // Tokyo
		{Lat: -33.8688, Lon: 151.2093}, // Sydney
	}

	for i := 0; i < len(coords); i++ {
		for j := i + 1; j < len(coords); j++ {
			hab := Haversine(coords[i], coords[j])
			hba := Haversine(coords[j], coords[i])

			if !almostEqual(hab, hba) {
				t.Errorf("Haversine not symmetric: d(%v,%v)=%v, d(%v,%v)=%v",
					i, j, hab, j, i, hba)
			}

			gab := GreatCircle(coords[i], coords[j])
			gba := GreatCircle(coords[j], coords[i])

			if !almostEqual(gab, gba) {
				t.Errorf("GreatCircle not symmetric: d(%v,%v)=%v, d(%v,%v)=%v",
					i, j, gab, j, i, gba)
			}

			vab, _ := VincentyKm(coords[i], coords[j])
			vba, _ := VincentyKm(coords[j], coords[i])

			if !almostEqualTolerance(vab, vba, 0.01) {
				t.Errorf("Vincenty not symmetric: d(%v,%v)=%v, d(%v,%v)=%v",
					i, j, vab, j, i, vba)
			}
		}
	}
}

// Benchmarks
func BenchmarkHaversine(b *testing.B) {
	nyc := Coord{Lat: 40.7128, Lon: -74.0060}
	london := Coord{Lat: 51.5074, Lon: -0.1278}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Haversine(nyc, london)
	}
}

func BenchmarkGreatCircle(b *testing.B) {
	nyc := Coord{Lat: 40.7128, Lon: -74.0060}
	london := Coord{Lat: 51.5074, Lon: -0.1278}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GreatCircle(nyc, london)
	}
}

func BenchmarkEquirectangular(b *testing.B) {
	coord1 := Coord{Lat: 40.0, Lon: -74.0}
	coord2 := Coord{Lat: 40.1, Lon: -74.1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Equirectangular(coord1, coord2)
	}
}

func BenchmarkVincenty(b *testing.B) {
	nyc := Coord{Lat: 40.7128, Lon: -74.0060}
	london := Coord{Lat: 51.5074, Lon: -0.1278}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Vincenty(nyc, london)
	}
}
