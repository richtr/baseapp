package models

// Datastore objects

type Like struct {
	LikeId          int
	PostId	        int
	UserId          int
}

type Follower struct {
	FollowerId      int
	UserId	        int
	FollowUserId    int
}

// JSON response objects

type SimpleJSONResponse struct {
  Status      string
  Message     string
}
