package appnexus

import (
	"errors"
	"fmt"
	"net/http"
)

// SegmentService handles all requests to the segment service API
type SegmentService struct {
	*Response
	client *Client
}

// Segment is an audience segment within the AppNexus console
type Segment struct {
	ID              int    `json:"id,omitempty"`
	Active          bool   `json:"active,omitempty"`
	Code            string `json:"code,omitempty"`
	State           string `json:"state,omitempty"`
	ShortName       string `json:"short_name"`
	Description     string `json:"description,omitempty"`
	MemberID        int    `json:"member_id"`
	Category        string `json:"category,omitempty"`
	ExpireMinutes   int    `json:"expire_minutes,omitempty"`
	AdvertiserID    int    `json:"advertiser_id,omitempty"`
	LastModified    string `json:"last_modified,omitempty"`
	Provider        string `json:"provider,omitempty"`
	ParentSegmentID int    `json:"parent_segment_id,omitempty"`
}

type segmentResponse struct {
	*http.Response
	Obj struct {
		Segment  `json:"segment,omitempty"`
		Segments []Segment `json:"segments,omitempty"`
		Error    string    `json:"error"`
		Status   string    `json:"status"`
		Service  string    `json:"service"`
		Rate     Rate      `json:"dbg_info"`
	} `json:"response"`
}

// Get a segment from the segment service by Member ID and Segment ID
func (s *SegmentService) Get(memberID int, segmentID int) (*Segment, error) {

	path := fmt.Sprintf("segment/%d?id=%d", memberID, segmentID)
	req, err := s.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	r := &segmentResponse{}
	_, err = s.client.do(req, r)
	if err != nil {
		return nil, err
	}

	segment := &r.Obj.Segment
	return segment, nil
}

// List available segments from your AppNexus console
func (s *SegmentService) List(memberID int, opt *ListOptions) ([]Segment, *Response, error) {
	u, err := addOptions(fmt.Sprintf("segment/%d", memberID), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.newRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	segments := &segmentResponse{}
	resp, err := s.client.do(req, segments)
	if err != nil {
		return nil, resp, err
	}

	return segments.Obj.Segments, resp, err
}

// Add a new segment
func (s *SegmentService) Add(memberID int, item *Segment) (*Response, error) {

	data := struct {
		Segment `json:"segment"`
	}{*item}

	req, err := s.client.newRequest("POST", fmt.Sprintf("segment/%d", memberID), data)

	if err != nil {
		return nil, err
	}

	result := &Response{}
	resp, err := s.client.do(req, result)
	if err != nil {
		return resp, err
	}

	item.ID = result.Obj.ID
	return result, nil
}

// Update an existing segment with new data
func (s *SegmentService) Update(memberID int, item Segment) (*Response, error) {

	data := struct {
		Segment `json:"segment"`
	}{item}

	if item.ID < 1 {
		return nil, errors.New("Update Segment requires a segment to have an ID already")
	}

	req, err := s.client.newRequest("PUT", fmt.Sprintf("segment/%d?id=%d", memberID, item.ID), data)

	if err != nil {
		return nil, err
	}

	result := &Response{}
	resp, err := s.client.do(req, result)
	if err != nil {
		return resp, err
	}

	return result, nil
}

// Delete the specified segment
func (s *SegmentService) Delete(memberID int, item Segment) error {

	data := struct {
		Segment `json:"segment"`
	}{item}

	if item.ID < 1 {
		return errors.New("Delete Segment requires a segment to have an ID already")
	}

	req, err := s.client.newRequest("DELETE", fmt.Sprintf("segment/%d", memberID), data)
	if err != nil {
		return err
	}

	_, err = s.client.do(req, nil)
	return err
}
