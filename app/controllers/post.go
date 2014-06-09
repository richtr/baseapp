package controllers

import (
	r "github.com/revel/revel"
	"github.com/richtr/baseapp/app/models"
	"github.com/richtr/baseapp/app/routes"
	markdown "github.com/russross/blackfriday"
	"html/template"
	"time"
)

type Post struct {
	Profile
}

func (c Post) loadPostById(id int) *models.Post {
	p, err := c.Txn.Get(models.Post{}, id)

	if err != nil || p == nil {
		return nil
	}

	return p.(*models.Post)
}

/*func (c Post) getConnectedPostOwner(post *models.Post) *models.Profile {
	profile := c.connected();
	if profile == nil {
		return nil
	}

	if profile.ProfileId != post.ProfileId {
		return nil
	}

	return profile
}*/

func (c Post) Show(username string, id int) r.Result {
	post := c.loadPostById(id)

	if post == nil {
		return c.NotFound("Post does not exist")
	}

	p, err := c.Txn.Get(models.Profile{}, post.ProfileId)
	if err != nil || p == nil {
		return c.NotFound("Profile does not exist")
	}
	profile := p.(*models.Profile)

	if profile == nil {
		return c.NotFound("Profile does not exist")
	}

	if profile.UserName != username {
		return c.NotFound("Post does not exist")
	}

	// Convert Content-field Markdown for rendering
	contentStr := markdown.MarkdownCommon(post.Content)
	contentStr = models.FormatContentMentions(contentStr)
	postContentHTML := template.HTML(contentStr)

	appName := r.Config.StringDefault("app.name", "BaseApp")

	title := post.Title + " / " + profile.Name + " on " + appName

	isOwner := false
	if user := c.connected(); user != nil && user.UserId == profile.User.UserId {
		isOwner = true
	}

	return c.Render(title, profile, post, postContentHTML, isOwner)
}

func (c Post) Create(username string) r.Result {
	profile := c.connected()
	if profile == nil || profile.UserName != username {
		c.Flash.Error("You must log in to access your account")
		return c.Redirect(routes.Account.Logout())
	}

	return c.Render(profile)
}

func (c Post) Save(username string, post models.Post) r.Result {
	profile := c.connected()
	if profile == nil || profile.UserName != username {
		c.Flash.Error("You must log in to access your account")
		return c.Redirect(routes.Account.Logout())
	}

	// Post component validation
	models.ValidatePostTitle(c.Validation, post.Title).Key("post.Title")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Could not create post")
		return c.Redirect(routes.Post.Create(username))
	}

	post.ProfileId = profile.ProfileId
	post.DateObj = time.Now()
	post.Status = "public"

	err := c.Txn.Insert(&post)
	if err != nil {
		panic(err)
	}

	c.Flash.Success("Post created")
	return c.Redirect(routes.Post.Show(username, post.PostId))
}

func (c Post) Edit(username string, id int) r.Result {
	post := c.loadPostById(id)
	if post == nil {
		return c.NotFound("Post does not exist")
	}

	profile := c.connected()
	if profile == nil || profile.UserName != username {
		c.Flash.Error("You must log in to access your account")
		return c.Redirect(routes.Account.Logout())
	}

	return c.Render(profile, post)
}

func (c Post) Update(username string, id int, post models.Post) r.Result {
	existingPost := c.loadPostById(id)
	if existingPost == nil {
		return c.NotFound("Post does not exist")
	}

	profile := c.connected()
	if profile == nil || profile.UserName != username {
		c.Flash.Error("You must log in to access your account")
		return c.Redirect(routes.Account.Logout())
	}

	// Post component validation
	models.ValidatePostTitle(c.Validation, post.Title).Key("post.Title")

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Could not update post")
		return c.Redirect(routes.Post.Edit(username, id))
	}

	// Update fields
	existingPost.Title = post.Title
	existingPost.ContentStr = post.ContentStr

	_, err := c.Txn.Update(existingPost)
	if err != nil {
		panic(err)
	}

	c.Flash.Success("Post updated")
	return c.Redirect(routes.Post.Show(username, existingPost.PostId))
}

func (c Post) Remove(username string, id int) r.Result {
	post := c.loadPostById(id)
	if post == nil {
		return c.NotFound("Post does not exist")
	}

	profile := c.connected()
	if profile == nil || profile.UserName != username {
		c.Flash.Error("You must log in to access your account")
		return c.Redirect(routes.Account.Logout())
	}

	return c.Render(profile, post)
}

func (c Post) Delete(username string, id int) r.Result {
	existingPost := c.loadPostById(id)
	if existingPost == nil {
		return c.NotFound("Post does not exist")
	}

	profile := c.connected()
	if profile == nil || profile.UserName != username {
		c.Flash.Error("You must log in to access your account")
		return c.Redirect(routes.Account.Logout())
	}

	_, err := c.Txn.Delete(existingPost)
	if err != nil {
		panic(err)
	}

	c.Flash.Success("Post removed")
	return c.Redirect(routes.Profile.Show(profile.UserName))
}

func (c Post) Like(username string, id int) r.Result {
	likeResponse := models.SimpleJSONResponse{"fail", ""}

	profile := c.connected()
	if profile == nil {
		likeResponse.Message = "You must log in to like a post"
		return c.RenderJson(likeResponse)
	}

	post := c.loadPostById(id)
	if post == nil {
		likeResponse.Message = "Post does not exist"
		return c.Render(likeResponse)
	}

	existingLike := &models.Like{}
	err := c.Txn.SelectOne(existingLike, `select LikeId from Like where UserId = ? and PostId = ?`, profile.User.UserId, post.PostId)

	if err == nil {

		likeResponse.Message = "You have already liked this post"

	} else {
		// Add like

		existingLike.LikeId = 0
		existingLike.UserId = profile.User.UserId
		existingLike.PostId = post.PostId

		lErr := c.Txn.Insert(existingLike)
		if lErr != nil {
			panic(lErr)
		}

		// Update like increment count on Post
		post.AggregateLikes += 1

		_, pErr := c.Txn.Update(post)
		if pErr != nil {
			panic(pErr)
		}

		likeResponse.Message = "You have now liked this post"
		likeResponse.Status = "success"
	}

	return c.RenderJson(likeResponse)

}
