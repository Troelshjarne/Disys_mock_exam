package main

import (
	mockPackage "github.com/Troelshjarne/Disys_mock_exam/increment"
)

type Server struct {
	mockPackage.UnimplementedCommunicationServer
}
