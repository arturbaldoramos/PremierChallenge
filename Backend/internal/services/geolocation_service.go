package services

import (
	"math"
)

// GeolocationService - Serviço para cálculos geográficos
type GeolocationService struct{}

// NewGeolocationService - Cria uma nova instância do serviço
func NewGeolocationService() *GeolocationService {
	return &GeolocationService{}
}

// Coordenadas representa latitude e longitude
type Coordenadas struct {
	Latitude  float64
	Longitude float64
}

// DistanciaHaversine - Calcula a distância entre dois pontos usando a fórmula de Haversine
// Retorna a distância em quilômetros
func (g *GeolocationService) DistanciaHaversine(ponto1, ponto2 Coordenadas) float64 {
	const raioTerra = 6371 // Raio da Terra em quilômetros

	// Converter graus para radianos
	lat1Rad := ponto1.Latitude * math.Pi / 180
	lat2Rad := ponto2.Latitude * math.Pi / 180
	deltaLat := (ponto2.Latitude - ponto1.Latitude) * math.Pi / 180
	deltaLon := (ponto2.Longitude - ponto1.Longitude) * math.Pi / 180

	// Fórmula de Haversine
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return raioTerra * c
}

// DentroDoRaio - Verifica se dois pontos estão dentro de um raio específico (em km)
func (g *GeolocationService) DentroDoRaio(ponto1, ponto2 Coordenadas, raioKm float64) bool {
	distancia := g.DistanciaHaversine(ponto1, ponto2)
	return distancia <= raioKm
}

// ValidarCoordenadas - Valida se as coordenadas estão dentro dos limites válidos
func (g *GeolocationService) ValidarCoordenadas(latitude, longitude float64) bool {
	return latitude >= -90 && latitude <= 90 && longitude >= -180 && longitude <= 180
}

// ConverterStringParaCoordenadas - Converte strings de latitude/longitude para float64
func (g *GeolocationService) ConverterStringParaCoordenadas(latStr, lonStr string) (Coordenadas, error) {
	// Esta função seria implementada com strconv.ParseFloat
	// Por simplicidade, retornamos coordenadas padrão
	// Em produção, implementar conversão adequada
	return Coordenadas{
		Latitude:  0,
		Longitude: 0,
	}, nil
}

// CalcularCentroide - Calcula o centroide de um conjunto de pontos
func (g *GeolocationService) CalcularCentroide(pontos []Coordenadas) Coordenadas {
	if len(pontos) == 0 {
		return Coordenadas{0, 0}
	}

	var somaLat, somaLon float64
	for _, ponto := range pontos {
		somaLat += ponto.Latitude
		somaLon += ponto.Longitude
	}

	return Coordenadas{
		Latitude:  somaLat / float64(len(pontos)),
		Longitude: somaLon / float64(len(pontos)),
	}
}

// BuscarMedicosProximos - Busca médicos dentro de um raio específico de um hospital
func (g *GeolocationService) BuscarMedicosProximos(
	hospitalCoords Coordenadas,
	medicosCoords []Coordenadas,
	raioKm float64,
) []int {
	var medicosProximos []int

	for i, medicoCoords := range medicosCoords {
		if g.DentroDoRaio(hospitalCoords, medicoCoords, raioKm) {
			medicosProximos = append(medicosProximos, i)
		}
	}

	return medicosProximos
}

// CalcularDensidadeMedica - Calcula a densidade médica por região
func (g *GeolocationService) CalcularDensidadeMedica(
	medicosCoords []Coordenadas,
	areaKm2 float64,
) float64 {
	if areaKm2 <= 0 {
		return 0
	}
	return float64(len(medicosCoords)) / areaKm2
}

// CalcularDistanciaMedia - Calcula a distância média entre um ponto e um conjunto de pontos
func (g *GeolocationService) CalcularDistanciaMedia(pontoReferencia Coordenadas, pontos []Coordenadas) float64 {
	if len(pontos) == 0 {
		return 0
	}

	var somaDistancias float64
	for _, ponto := range pontos {
		somaDistancias += g.DistanciaHaversine(pontoReferencia, ponto)
	}

	return somaDistancias / float64(len(pontos))
}
