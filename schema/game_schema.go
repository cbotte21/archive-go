package schema

// Game struct
type Game struct { //Payload
	Id        string `bson:"_id,omitempty" json:"_id,omitempty" redis:"_id"`
	Viktor    string `bson:"viktor,omitempty" json:"viktor,omitempty" redis:"viktor"`
	Opponent  string `bson:"opponent,omitempty" json:"opponent,omitempty" redis:"opponent"`
	DeltaElo  int32  `bson:"deltaElo,omitempty" json:"deltaElo,omitempty" redis:"deltaElo"`
	Timestamp int64  `bson:"timestamp,omitempty" json:"timestamp,omitempty" redis:"timestamp"`
}

func (game Game) Database() string {
	return ""
}

func (game Game) Collection() string {
	return ""
}

func (game Game) Key() string {
	return game.Id
}
