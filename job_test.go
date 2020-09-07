package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

type mockJobInterface struct {
	list               *v1.JobList
	listErr, deleteErr error
}

func (m mockJobInterface) Create(ctx context.Context, job *v1.Job, opts metav1.CreateOptions) (*v1.Job, error) {
	return nil, nil
}
func (m mockJobInterface) Update(ctx context.Context, job *v1.Job, opts metav1.UpdateOptions) (*v1.Job, error) {
	return nil, nil
}
func (m mockJobInterface) UpdateStatus(ctx context.Context, job *v1.Job, opts metav1.UpdateOptions) (*v1.Job, error) {
	return nil, nil
}
func (m mockJobInterface) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return m.deleteErr
}
func (m mockJobInterface) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return nil
}
func (m mockJobInterface) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Job, error) {
	return nil, nil
}
func (m mockJobInterface) List(ctx context.Context, opts metav1.ListOptions) (*v1.JobList, error) {
	return m.list, m.listErr
}
func (m mockJobInterface) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return nil, nil
}
func (m mockJobInterface) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Job, err error) {
	return nil, nil
}

func Test_JobCleaner_Delete(t *testing.T) {
	ttl := 1 * time.Second
	t.Run("job does not have conditions", func(t *testing.T) {
		jobList := &v1.JobList{
			Items: []v1.Job{
				v1.Job{
					Status: v1.JobStatus{
						Conditions: []v1.JobCondition{},
					},
				},
			},
		}
		jc := NewJobCleaner(mockJobInterface{list: jobList, listErr: nil, deleteErr: nil}, ttl)
		require.NoError(t, jc.Delete(context.Background()))
	})
	t.Run("job list returns error", func(t *testing.T) {
		jobList := &v1.JobList{}
		jc := NewJobCleaner(mockJobInterface{list: jobList, listErr: fmt.Errorf("list error"), deleteErr: nil}, ttl)
		require.Error(t, jc.Delete(context.Background()))
	})
	t.Run("job has conditions", func(t *testing.T) {
		t.Run("passed ttl", func(t *testing.T) {
			jobList := &v1.JobList{
				Items: []v1.Job{
					v1.Job{
						Status: v1.JobStatus{
							Conditions: []v1.JobCondition{v1.JobCondition{
								Type:          v1.JobComplete,
								LastProbeTime: metav1.Time{Time: time.Now().Add(-10 * time.Second)},
							}},
						},
					},
				},
			}
			t.Run("job delete success", func(t *testing.T) {
				jc := NewJobCleaner(mockJobInterface{list: jobList, listErr: nil, deleteErr: nil}, ttl)
				require.NoError(t, jc.Delete(context.Background()))
			})
			t.Run("job delete returns error", func(t *testing.T) {
				jc := NewJobCleaner(mockJobInterface{list: jobList, listErr: nil, deleteErr: fmt.Errorf("delete error")}, ttl)
				require.Error(t, jc.Delete(context.Background()))
			})
		})
		t.Run("does not passed ttl", func(t *testing.T) {
			jobList := &v1.JobList{
				Items: []v1.Job{
					v1.Job{
						Status: v1.JobStatus{
							Conditions: []v1.JobCondition{v1.JobCondition{
								Type:          v1.JobComplete,
								LastProbeTime: metav1.Time{Time: time.Now()},
							}},
						},
					},
				},
			}
			jc := NewJobCleaner(mockJobInterface{list: jobList, listErr: nil, deleteErr: fmt.Errorf("delete error")}, ttl)
			// does not execute Delete
			require.NoError(t, jc.Delete(context.Background()))
		})
	})
}
