package internal

import (
	"context"
	"github.com/cbotte21/archive-go/pb"
	"github.com/cbotte21/archive-go/schema"
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Archive struct {
	Games          *datastore.MongoClient[schema.Game]
	ServiceRecords *datastore.MongoClient[schema.SVCRecord]
	pb.UnimplementedArchiveServiceServer
}

func (archive *Archive) mustEmbedUnimplementedArchiveServiceServer() {
	//TODO implement me
	panic("implement me")
}

func NewArchive(games *datastore.MongoClient[schema.Game], svcRecords *datastore.MongoClient[schema.SVCRecord]) Archive {
	return Archive{Games: games, ServiceRecords: svcRecords}
}

func (archive *Archive) History(context context.Context, player *pb.Player) (*pb.PriorMatches, error) {
	history, err := archive.ServiceRecords.Find(schema.SVCRecord{Player: player.GetXId()})
	if err != nil {
		return nil, err
	}
	return &pb.PriorMatches{Id: history.Games}, nil
}

func (archive *Archive) Match(context context.Context, matchId *pb.MatchId) (*pb.MatchInfo, error) {
	game, err := archive.Games.Find(schema.Game{Id: matchId.GetId()})
	if err != nil {
		return nil, err
	}
	return &pb.MatchInfo{
		Viktor:    &pb.Player{XId: game.Viktor},
		Opponent:  &pb.Player{XId: game.Opponent},
		DeltaElo:  game.DeltaElo,
		Timestamp: &timestamppb.Timestamp{Seconds: game.Timestamp},
	}, nil
}

func (archive *Archive) Record(context context.Context, matchInfo *pb.MatchInfo) (*pb.Void, error) {
	candideGame := schema.Game{
		Viktor:    matchInfo.GetViktor().GetXId(),
		Opponent:  matchInfo.GetOpponent().GetXId(),
		DeltaElo:  matchInfo.GetDeltaElo(),
		Timestamp: matchInfo.GetTimestamp().Seconds,
	}

	err := archive.Games.Create(candideGame)
	if err != nil {
		return nil, err
	}

	game, err := archive.Games.Find(candideGame) // To get _id
	if err != nil {
		return nil, err
	}

	viktorHistory, err := archive.ServiceRecords.Find(schema.SVCRecord{Player: matchInfo.GetViktor().GetXId()})
	opponentHistory, err := archive.ServiceRecords.Find(schema.SVCRecord{Player: matchInfo.GetOpponent().GetXId()})

	viktorHistory.Games = append(viktorHistory.Games, game.Key())
	opponentHistory.Games = append(opponentHistory.Games, game.Key())

	_ = archive.ServiceRecords.Update(schema.SVCRecord{Player: viktorHistory.Key()}, viktorHistory)
	_ = archive.ServiceRecords.Update(schema.SVCRecord{Player: opponentHistory.Key()}, opponentHistory)
	return nil, nil
}
