package services

// type CommentService struct {
// 	Repo *repositories.CommentRepository
// }

// func NewCommentService(repo repositories.CommentRepository) *CommentService {
// 	return &CommentService{
// 		Repo: re,
// 	}
// }

// func (cs *CommentService) GetAllComments() ([]map[string]interface{}, error) {
// 	return cs.Repo.FetchComments()
// }

// func (cs *CommentService) AddComment(fileID, userID, content string) error {
// 	// Business logic, validation, etc.
// 	if fileID == "" || userID == "" || content == "" {
// 		return &InvalidInputError{"All fields are required"}
// 	}
// 	return cs.Repo.InsertComment(fileID, userID, content)
// }

// // Custom error handling
// type InvalidInputError struct {
// 	Message string
// }

// func (e *InvalidInputError) Error() string {
// 	return e.Message
// }
