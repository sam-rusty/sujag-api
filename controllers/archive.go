package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/user/sujag/util"
	"github.com/user/sujag/models"
)

func (c *Controllers) getArchivePosts(ctx *gin.Context) {

	var Posts []models.Post

	SelectField := bson.M{"title": 1, "thumbnail": 1, "excerpt": 1, "district": 1, "slug": 1, "published_on": 1}

	where := bson.M{}

	Params := util.GetParams(ctx)

	Page := 0.0

	if Params["page"] != nil {
		Page = Params["page"].(float64)
	}

	if Params["type"] != nil {
		where["type"] = Params["type"].(string)
	} else {
		where["type"] = "feature"
	}

	if Params["category"] != nil {
		where["category"] = Params["category"].(string)
	}

	if Params["district"] != nil && Params["district"] != "all" {
		where["district"] = Params["district"].(string)
	}

	if Params["author"] != nil {
		where["author"] = Params["author"].(string)
	}

	Limit := 20
	Skip := 0

	if Page != 0 {
		Skip = int(Page) * 20
	}

	where["status"] = 1

	c.App.DB.C("posts").Find(where).Select(SelectField).Sort("-published_on").Skip(Skip).Limit(Limit).All(&Posts)

	Total, _ := c.App.DB.C("posts").Find(where).Count()

	ctx.JSON(http.StatusOK, bson.M{
		"posts": Posts,
		"total": Total,
	})

}

func (c *Controllers) SearchPosts(ctx *gin.Context) {
	Params := util.GetParams(ctx)

	Page := 0.00
	if Params["page"] != nil {
		Page = Params["page"].(float64)
	}
	var Posts models.Post
	where := bson.M{}
	where["status"] = 1
	Limit := 20
	Skip := int(Page) * 20
	SearchKeywords := Params["search_keywords"].(string)
	if SearchKeywords != "" {
		where["$text"] = bson.M{"$search": SearchKeywords}
	}
	c.App.DB.C("posts").Find(where).Skip(Skip).Limit(Limit).All(&Posts)
	Total, _ := c.App.DB.C("posts").Find(where).Sort("-published_on").Skip(Skip).Limit(Limit).Count()
	ctx.JSON(http.StatusOK, bson.M{
		"posts": Posts,
		"total": Total,
	})
}
