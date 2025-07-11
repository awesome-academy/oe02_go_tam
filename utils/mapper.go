package utils

import (
	"oe02_go_tam/models"
	"oe02_go_tam/responses"
)

func MapReviewToResponse(r models.Review) responses.ReviewResponse {
	// Map comments
	var commentRes []responses.CommentResponse
	for _, c := range r.Comments {
		commentRes = append(commentRes, responses.CommentResponse{
			ID: c.ID,
			User: responses.UserMiniResponse{
				ID:    c.User.ID,
				Name:  c.User.Name,
				Email: c.User.Email,
			},
			Content: c.Content,
		})
	}

	// Map likes
	var likeRes []responses.LikeResponse
	for _, l := range r.Likes {
		likeRes = append(likeRes, responses.LikeResponse{
			ID: l.ID,
			User: responses.UserMiniResponse{
				ID:    l.User.ID,
				Name:  l.User.Name,
				Email: l.User.Email,
			},
		})
	}

	// Map review
	return responses.ReviewResponse{
		ID:      r.ID,
		Rating:  r.Rating,
		Content: r.Content,
		User: responses.UserMiniResponse{
			ID:    r.User.ID,
			Name:  r.User.Name,
			Email: r.User.Email,
		},
		Comments: commentRes,
		Likes:    likeRes,
	}
}
