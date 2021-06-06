package myMongo

import (
	"armani_follow/pkg/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"time"
)

type Followers struct {
	MyClient 		*mongo.Client
}

func (f *Followers) Follow(me models.User, someone models.User) error {

	collection := f.MyClient.Database("follow").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err := collection.FindOne(ctx, bson.M{"_id": me.ID})
	if err.Err() == mongo.ErrNoDocuments{
		var response1 models.MongoResponse
		res1 := models.MongoResponse{
			IDUser:      me.ID,
			Subscribers: []models.User{},
			Subscribing: []models.User{someone},
		}
		collection.InsertOne(ctx,res1) // works fine
		err1 := collection.FindOne(ctx,  bson.M{"_id": someone.ID})
		fmt.Println(someone.ID)
		if err1.Err() == mongo.ErrNoDocuments{
			res2 := models.MongoResponse{
				IDUser:      someone.ID,
				Subscribers: []models.User{me},
				Subscribing: []models.User{},
			}
			collection.InsertOne(ctx,res2)
		}else{

			err2 := err1.Decode(&response1)
			if err2 != nil{
				return err2
			}
			response1.Subscribers = append(response1.Subscribers, me)


			update := bson.M{"$set": bson.M{"subscribers": response1.Subscribers}}
			result, err3 := collection.UpdateOne(
				ctx,
				bson.M{"_id": someone.ID},
				update,
			)
			fmt.Println("update when my lists not exist and someones has" ,result.MatchedCount, result.ModifiedCount)
			if err3 != nil{
				return err3
			}

			return nil
		}
		return nil
	}else{
		var response models.MongoResponse
		err2 := err.Decode(&response)
		if err2 != nil{
			return err2
		}
		response.Subscribing = append(response.Subscribing, someone)


		//update := bson.NewRegistryBuilder(
		//	bson..SubDocumentFromElements(
		//		"$set",
		//		bson.EC.Double("pi", 3.14159),
		//	),
		//)

		update := bson.M{"$set": bson.M{"subscribing": response.Subscribing}}

		_, err3 := collection.UpdateOne(
			ctx,
			bson.M{"_id": me.ID},
			update,
		)
		//fmt.Println("update when my lists not exist and someones has", result.MatchedCount, result.ModifiedCount)
		if err3 != nil {
			return err3
		}
		err1 := collection.FindOne(ctx, bson.M{"_id": someone.ID})
		if err1.Err() == mongo.ErrNoDocuments{
			res2 := models.MongoResponse{
				IDUser:      someone.ID,
				Subscribers: []models.User{me},
				Subscribing: []models.User{},
			}
			collection.InsertOne(ctx,res2)

		}else{

			err2 := err1.Decode(&response)
			if err2 != nil{
				return err2
			}
			response.Subscribers = append(response.Subscribers, me)

			update := bson.M{"$set": bson.M{"subscribers": response.Subscribers}}
			result, err3 := collection.UpdateOne(
				ctx,
				bson.M{"_id": someone.ID},
				update,
			)
			fmt.Println("update when my lists exists and someones also" ,result.MatchedCount, result.ModifiedCount)
			if err3 != nil{
				return err3
			}



			return nil
		}


		return nil
	}

}


func (f *Followers) GetFollowers(me models.User) ([]models.User,error) {
	collection := f.MyClient.Database("follow").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err := collection.FindOne(ctx, bson.M{"_id": me.ID})
	if err.Err() == mongo.ErrNoDocuments{
		return nil, errors.New("No one is following you")
	}else{

		var response models.MongoResponse
		err2 := err.Decode(&response)
		if err2 != nil{
			return nil, err2
		}
		return response.Subscribing, nil

	}
}




func (f *Followers) Unfollow(me models.User, someone models.User) error {
	collection := f.MyClient.Database("follow").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	err := collection.FindOne(ctx, bson.M{"_id": me.ID})

	if err.Err() == mongo.ErrNoDocuments{
		return err.Err()
	}else{

		var response models.MongoResponse
		err2 := err.Decode(&response)
		if err2 != nil{
			return err2
		}





	}

	return nil
}
