package taskstore

import (
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	// Creates a new taskstore and populates it with two test tasks
	ts := New()
	id_0 := ts.CreateTask(
		"Test task 1",
		[]string{"tag 1.1", "tag 1.2"},
		time.Date(2023, time.September, 5, 0, 0, 0, 0, time.UTC),
	)
	id_1 := ts.CreateTask(
		"Test task 2",
		[]string{"tag 2.1", "tag 2.2"},
		time.Date(2023, time.September, 5, 0, 0, 0, 0, time.UTC),
	)

	// Checks the size of the taskstore
	if len(ts.tasks) != 2 {
		t.Errorf("server returned the wrong length: got %v wanted '2'", len(ts.tasks))
	}

	// Checks if the returned ID mathces the expected values
	if id_0 != 0 {
		t.Errorf("server returned the wrong task ID: got %v wanted '0'", id_0)
	}
	if id_1 != 1 {
		t.Errorf("server returned the wrong task ID: got %v wanted '1'", id_0)
	}

	// Gets the two test tasks from the taskstore
	task_1, err := ts.GetTask(id_0)
	if err != nil {
		t.Errorf("server returned an error: got %v wanted 'NIL'", err)
	}

	// Compares each field of the test tasks
	if task_1.Text != "Test task 1" {
		t.Errorf("server returned the wrong test text: got %v wanted 'Test task 1'", task_1.Text)
	}
	if task_1.Tags[0] != "tag 1.1" {
		t.Errorf("server returned the wrong test tag: got %v wanted 'Test task 1'", task_1.Tags[0])
	}
	if task_1.Due != time.Date(2023, time.September, 5, 0, 0, 0, 0, time.UTC) {
		t.Errorf("server returned the wrong test date: got %v wanted '2023/09/05'", task_1.Due)
	}

	// Gets all tasks
	allTasks := ts.GetAllTasks()
	if len(allTasks) != 2 {
		t.Errorf("server returned the wrong number of tasks: got %v wanted 2", len(allTasks))
	}

	// Gets an unexisting task
	_, err = ts.GetTask(2)
	if err == nil {
		t.Errorf("server returned the wrong error: got 'nil' wanted 'runtime error'")
	}
}

func TestDeleteTask(t *testing.T) {
	// Creates a new taskstore and populates it with two test tasks
	ts := New()
	id_0 := ts.CreateTask(
		"Test task 1",
		[]string{"tag 1.1", "tag 1.2"},
		time.Date(2023, time.September, 5, 0, 0, 0, 0, time.UTC),
	)
	_ = ts.CreateTask(
		"Test task 2",
		[]string{"tag 2.1", "tag 2.2"},
		time.Date(2023, time.September, 5, 0, 0, 0, 0, time.UTC),
	)

	// Deletes an unexisting task
	err := ts.DeleteTask(2)
	if err == nil {
		t.Errorf("server returned the wrong error: got 'NIL' wanted 'ERROR'")
	}

	// Deletes an existing task
	err = ts.DeleteTask(id_0)
	if len(ts.tasks) != 1 {
		t.Errorf("server returned the wrong number of tasks: got %v wanted '1'", len(ts.tasks))
	}

	// Deletes all tasks
	err = ts.DeleteAllTasks()
	if len(ts.tasks) != 0 {
		t.Errorf("server returned the wrong number of tasks: got %v wanted '0'", len(ts.tasks))
	}

}
