package main

import (
	"fmt"
	"os"

	"github.com/adwww/appnexus"
	"github.com/fatih/color"
)

func main() {

	// Create a new client and login to authentication service;
	c := AppNexus.NewClient(nil)
	err := c.Login("<<<username@example.com>>>", "<<<password>>>")
	if err != nil {
		color.Red(err.Error())
		os.Exit(2)
	}

	// Load the default AppNexus member for this login:
	member, err := c.Members.GetDefault()
	if err != nil {
		color.Red(err.Error())
		os.Exit(3)
	}

	fmt.Println("\n\nConnected as member", color.YellowString(member.Name))

	// List the first 20 segments available to this member:
	segments, _, err := c.Segments.List(member.ID, &AppNexus.ListOptions{StartElement: 0, NumElements: 20})
	if err != nil {
		color.Red(err.Error())
		os.Exit(4)
	}

	fmt.Printf("\n\nFound %d segments\n", len(segments))

	for _, s := range segments {
		fmt.Println(" - ", s.ID, color.YellowString(s.ShortName), s.Code)
	}

	// Create a new segment:
	newSegment := AppNexus.Segment{
		ShortName:   "Test Seggy",
		MemberID:    member.ID,
		Description: "Test segment from the Go AppNexus API",
		Active:      false,
	}

	fmt.Print("\n\nCreating a new ", color.YellowString(newSegment.ShortName), " segment")
	resp, err := c.Segments.Add(member.ID, &newSegment)
	if err != nil {
		color.Red(" FAILED\n" + err.Error())
		os.Exit(5)
	}

	color.Green(" OK!")
	fmt.Printf("Test segment with ID %s created succsesfully\n", color.YellowString(fmt.Sprintf("%d", resp.Obj.ID)))

	// Update our new segment:
	newSegment.Code = "go_client_test"
	fmt.Print("\n\nUpdating segment ", color.YellowString(fmt.Sprintf("%d", resp.Obj.ID)), " code")

	resp, err = c.Segments.Update(member.ID, newSegment)
	if err != nil {
		color.Red(" FAILED\n" + err.Error())
		os.Exit(6)
	}

	color.Green(" OK!")

	// Delete the test segment:
	fmt.Print("\n\nDelete the ", color.YellowString(newSegment.ShortName), " test segment")
	err = c.Segments.Delete(member.ID, newSegment)
	if err != nil {
		color.Red(" FAILED\n" + err.Error())
		os.Exit(7)
	}

	color.Green(" OK!")
}
