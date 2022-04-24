package service

import "GolangFrameDemo/internal/dao"

var service *StuService

type StuService struct {
	d dao.Dao
}

func New(d dao.Dao)(s *StuService,cfg func(),err error)  {
	s = &StuService{d:d}
	cfg = s.Close
	service = s
	return
}

func (s *StuService)Close()  {
	s.d.Close()
}

func (s *StuService)GetStudentById(id int) (*dao.Student,error) {
	return s.d.GetStudent(id)
}

func GetService() *StuService  {
	return service
}
