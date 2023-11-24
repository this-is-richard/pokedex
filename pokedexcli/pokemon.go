package pokedexcli

import (
	"math/rand"

	"github.com/this-is-richard/pokedex/pokeapi"
)

type Pokemon struct {
	pokeapi.Pokemon
}

func NewPokemon(pokemon pokeapi.Pokemon) Pokemon {
	return Pokemon{
		Pokemon: pokemon,
	}
}

func (p *Pokemon) Catch() bool {
	return rand.Float64()*100/float64(p.BaseExperience) > 0.5
}
