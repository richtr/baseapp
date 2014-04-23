package controllers

import (
	r "github.com/revel/revel"
	"github.com/richtr/baseapp/app/models"
  "github.com/richtr/baseapp/app/routes"
  "time"
  markdown "github.com/russross/blackfriday"
  "html/template"
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

func (c Post) getConnectedPostOwner(post *models.Post) *models.Profile {
  profile := c.connected();
  if profile == nil {
    return nil
  }

  if profile.ProfileId != post.ProfileId {
    return nil
  }

  return profile
}

func (c Post) Index() r.Result {
	return c.NotFound("Post does not exist", 404)
}

func (c Post) Show(id int) r.Result {
  post := c.loadPostById(id)

	if post == nil {
		return c.NotFound("Post does not exist")
	}

  profile := c.loadProfileById(post.ProfileId)

	if profile == nil {
		return c.NotFound("Profile does not exist")
	}

  if profile.ProfileId != post.ProfileId {
    return c.NotFound("Post does not exist")
  }

  // Convert Content-field Markdown for rendering
  contentStr := markdown.MarkdownCommon(post.Content)
  postContentHTML := template.HTML(contentStr)

  appName := r.Config.StringDefault("app.name", "BaseApp")

	title := post.Title + " / " + profile.Name + " on " + appName

  isOwner := false
  if user := c.connected(); user != nil && user.UserId == profile.User.UserId {
    isOwner = true
  }

	return c.Render(title, profile, post, postContentHTML, isOwner)
}

func (c Post) Create() r.Result {
  profile := c.connected();
  if profile == nil {
    c.Flash.Error("You must log in to access your account")
    return c.Redirect(routes.Account.Logout())
  }

	return c.Render(profile)
}

func (c Post) Save(post models.Post) r.Result {
  profile := c.connected();
  if profile == nil {
    c.Flash.Error("You must log in to access your account")
    return c.Redirect(routes.Account.Logout())
  }

  // Post component validation
  models.ValidatePostTitle(c.Validation, post.Title).Key("post.Title")

  if c.Validation.HasErrors() {
    c.Validation.Keep()
    c.FlashParams()
    c.Flash.Error("Could not create post")
    return c.Redirect(routes.Post.Create())
  }

  post.ProfileId = profile.ProfileId
  post.DateObj = time.Now()
  post.Status = "public"

  err := c.Txn.Insert(&post)
  if err != nil {
    panic(err)
  }

  c.Flash.Success("Post created")
  return c.Redirect(routes.Post.Show(post.PostId))
}

func (c Post) Edit(id int) r.Result {
  post := c.loadPostById(id)
	if post == nil {
		return c.NotFound("Post does not exist")
	}

  profile := c.getConnectedPostOwner(post)
  if profile == nil {
    c.Flash.Error("You must log in to access your account");
    return c.Redirect(routes.Account.Logout())
  }

  return c.Render(profile, post)
}

func (c Post) Update(id int, post models.Post) r.Result {
  existingPost := c.loadPostById(id)
	if existingPost == nil {
		return c.NotFound("Post does not exist")
	}

  profile := c.getConnectedPostOwner(existingPost)
  if profile == nil {
    c.Flash.Error("You must log in to access your account");
    return c.Redirect(routes.Account.Logout())
  }

  // Post component validation
  models.ValidatePostTitle(c.Validation, post.Title).Key("post.Title")

  if c.Validation.HasErrors() {
    c.Validation.Keep()
    c.FlashParams()
    c.Flash.Error("Could not update post")
    return c.Redirect(routes.Post.Edit(id))
  }

  // Update fields
  existingPost.Title = post.Title;
  existingPost.ContentStr = post.ContentStr;

  _, err := c.Txn.Update(existingPost)
  if err != nil {
    panic(err)
  }

  c.Flash.Success("Post updated")
  return c.Redirect(routes.Post.Show(existingPost.PostId))
}

func (c Post) Remove(id int) r.Result {
  post := c.loadPostById(id)
	if post == nil {
		return c.NotFound("Post does not exist")
	}

  profile := c.getConnectedPostOwner(post)
  if profile == nil {
    c.Flash.Error("You must log in to access your account");
    return c.Redirect(routes.Account.Logout())
  }

  return c.Render(profile, post)
}

func (c Post) Delete(id int) r.Result {
  existingPost := c.loadPostById(id)
	if existingPost == nil {
		return c.NotFound("Post does not exist")
	}

  profile := c.getConnectedPostOwner(existingPost)
  if profile == nil {
    c.Flash.Error("You must log in to access your account");
    return c.Redirect(routes.Account.Logout())
  }

  _, err := c.Txn.Delete(existingPost)
  if err != nil {
    panic(err)
  }

  c.Flash.Success("Post removed")
  return c.Redirect(routes.Profile.Show(profile.ProfileId))
}