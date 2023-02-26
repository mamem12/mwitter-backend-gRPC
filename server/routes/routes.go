package routes

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/mwitter-backend-gRPC/proto/v1/mweets"
	"github.com/mwitter-backend-gRPC/server/dblayer"
	"github.com/mwitter-backend-gRPC/server/models"
	"github.com/mwitter-backend-gRPC/server/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type MweetRoutes struct {
	pb.MweetsServer
}

var db dblayer.DBLayer

func GetServer() (*MweetRoutes, error) {

	initdb, err := dblayer.NewORM("test", gorm.Config{})

	if err != nil {
		return nil, err
	}

	db = initdb

	return &MweetRoutes{}, nil
}

func (mr *MweetRoutes) GetAllMweeter(req *pb.ListMweetsRequest, stream pb.Mweets_GetAllMweeterServer) error {

	mweets, err := db.GetAllMweeter()

	if err != nil {
		return status.Errorf(codes.Internal, err.Error())

	}

	for _, mweet := range mweets {
		id := strconv.Itoa(int(mweet.ID))
		stream.Send(&pb.Mweet{
			ID:      &id,
			Image:   &mweet.Image,
			Content: &mweet.Content,
			UserId:  &mweet.UserId,
		})
	}

	return nil
}

func (mr *MweetRoutes) CreateMweet(ctx context.Context, req *pb.CreateMweetRequest) (*pb.MweetResponse, error) {

	userId := req.GetUserId()

	cookieString := utils.GetCookie(ctx)
	// fmt.Println(cookieString)
	cookieUserId, err := utils.ParseCookie(cookieString, "ID")

	if err != nil || cookieUserId != userId {
		return nil, status.Error(codes.Unavailable, "not equal user")
	}

	existUserId, err := strconv.Atoi(cookieUserId)

	if err != nil {
		return nil, err
	}

	if !db.ExistUser(existUserId) {
		return nil, status.Error(codes.NotFound, "not found user")
	}

	image := req.GetImage()
	content := req.GetContent()

	if userId == "" || (content == "" && image == "") {
		return nil, status.Error(codes.Unavailable, "bad request")
	}

	mweet := &models.Mweet{
		Image:   image,
		Content: content,
		UserId:  userId,
	}

	err = db.CreateMweet(mweet)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	fmt.Printf("%+v\n", mweet)

	mweetId := strconv.Itoa(int(mweet.Model.ID))

	responseMweet := &pb.Mweet{
		ID:      &mweetId,
		Image:   &mweet.Image,
		Content: &mweet.Content,
		UserId:  &mweet.UserId,
	}

	return &pb.MweetResponse{
		Mweet: responseMweet,
	}, nil
}

func (mr *MweetRoutes) UpdateMweet(ctx context.Context, req *pb.UpdateMweetRequest) (*pb.MweetResponse, error) {

	userId := req.GetUserId()

	cookieString := utils.GetCookie(ctx)
	cookieUserId, err := utils.ParseCookie(cookieString, "ID")

	if err != nil {
		return nil, err
	}

	existUserId, err := strconv.Atoi(cookieUserId)

	if err != nil {
		return nil, err
	}

	if !db.ExistUser(existUserId) {
		return nil, status.Error(codes.NotFound, "not found user")
	}

	if userId == "" || cookieUserId != userId {
		return nil, status.Errorf(codes.Unavailable, "유효하지 않은 요청입니다.")
	}

	mweetId, err := strconv.Atoi(req.GetID())

	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}

	findMweet, err := db.GetMweeterById(mweetId)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if req.GetContent() != "" {
		findMweet.Content = req.GetContent()
	}

	err = db.UpdateMweet(mweetId, findMweet)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	responseId := strconv.Itoa(int(findMweet.Model.ID))

	return &pb.MweetResponse{
		Mweet: &pb.Mweet{
			ID:      &responseId,
			Image:   &findMweet.Image,
			Content: &findMweet.Content,
			UserId:  &userId,
		},
	}, nil
}

func (mr *MweetRoutes) DeleteMweet(ctx context.Context, req *pb.DeleteMweetRequest) (*pb.DeleteMweetResponse, error) {

	mweetId, err := strconv.Atoi(req.GetID())

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	cookieString := utils.GetCookie(ctx)
	cookieUserId, err := utils.ParseCookie(cookieString, "ID")

	if err != nil {
		return nil, err
	}

	existMweet, err := db.GetMweeterById(mweetId)

	if err != nil {
		return nil, err
	}

	if existMweet.UserId != cookieUserId {
		return nil, status.Error(codes.PermissionDenied, "유효하지 않은 권한입니다.")
	}

	err = db.DeleteMweet(mweetId)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.DeleteMweetResponse{
		Success: true,
	}, nil
}
