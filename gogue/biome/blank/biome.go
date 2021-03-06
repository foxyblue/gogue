package blank

import (
	area "github.com/foxyblue/gogue/gogue/biome"
	"github.com/foxyblue/gogue/gogue/biome/factory"
	"github.com/foxyblue/gogue/gogue/entity"
)

const (
	biomeName = "blank"

	defaultBiomeCreatures = 3
	defaultBiomeRooms     = 3
)

// BiomeParameters represents all the configuration options
// for the biome.
type BiomeParameters struct {
	biomeCreatures int
	biomeRooms     int
	start          *area.Coord
	end            *area.Coord
	x              int
	y              int
	maxX           int
	maxY           int
}

func init() {
	factory.Register(biomeName, &blankBiomeFactory{})
}

// blankBiomeFactory implements the factory.biomeFactory interface
type blankBiomeFactory struct{}

func (factory *blankBiomeFactory) Create(parameters map[string]interface{}) (area.Biome, error) {
	return fromParameters(parameters)
}

type biome struct {
	parameters   BiomeParameters
	Grid         area.Grid
	ListEntities []entity.Entity
}

func fromParameters(parameters map[string]interface{}) (area.Biome, error) {
	params, err := fromParametersImpl(parameters)
	if err != nil || params == nil {
		return nil, err
	}
	return New(*params), nil
}

func fromParametersImpl(parameters map[string]interface{}) (*BiomeParameters, error) {
	var (
		creatures = defaultBiomeCreatures
		rooms     = defaultBiomeRooms
		start     = &area.Coord{X: 10, Y: 10}
		end       = &area.Coord{X: 20, Y: 20}
	)

	if startXY, ok := parameters["start"]; ok {
		start = startXY.(*area.Coord)
	}

	params := &BiomeParameters{
		biomeCreatures: creatures,
		biomeRooms:     rooms,
		start:          start,
		end:            end,
		maxX:           parameters["maxX"].(int) - 1,
		maxY:           parameters["maxY"].(int) - 1,
		x:              parameters["x"].(int),
		y:              parameters["y"].(int),
	}
	return params, nil
}

// New returns a constructed biome, if the linter fails it means
// we haven't implemented all the required methods on the biome
func New(params BiomeParameters) area.Biome {
	w := params.maxX - params.x
	h := params.maxY - params.y
	grid := area.NewGrid(params.x, params.y, w, h)
	return &biome{
		parameters:   params,
		Grid:         *grid,
		ListEntities: make([]entity.Entity, 2), //params.biomeCreatures),
	}
}

func (b *biome) Generate() {
	randomCoord := area.RandomCoord(
		b.parameters.x, b.parameters.maxX,
		b.parameters.y, b.parameters.maxY)

	rabbit := NewRabbit(randomCoord.X, randomCoord.Y)
	b.ListEntities[0] = rabbit

	sword := NewSword(6, 6)
	b.ListEntities[1] = sword

	g := b.Grid
	room := []*area.Coord{
		{X: 5, Y: 5},
		{X: 4, Y: 5},
		{X: 3, Y: 5},
		{X: 2, Y: 5},
	}

	for x, row := range g.Tiles {
		for y := range row {
			if area.IsIn(x, y, room) {
				row[y] = area.WallTile(x, y)
			} else {
				row[y] = area.EmptyTile(x, y)
			}
		}
	}
}

func (b *biome) GetGrid() area.Grid {
	return b.Grid
}

func (b *biome) GetEntities() []entity.Entity {
	return b.ListEntities
}

func (b *biome) StartLocation() *area.Coord {
	return b.parameters.start
}

func (b *biome) EndLocation() *area.Coord {
	return b.parameters.end
}
