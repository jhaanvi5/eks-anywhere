package vsphere

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"

	"github.com/aws/eks-anywhere/pkg/config"
	"github.com/aws/eks-anywhere/pkg/govmomi"
	"github.com/aws/eks-anywhere/pkg/govmomi/mocks"
)

func TestValidatorValidatePrivs(t *testing.T) {
	v := Validator{}

	ctrl := gomock.NewController(t)
	vsc := mocks.NewMockVSphereClient(ctrl)

	ctx := context.Background()
	networkPath := "/Datacenter/network/path/foo"

	objects := []PrivAssociation{
		{
			objectType:   govmomi.VSphereTypeNetwork,
			privsContent: config.VSphereUserPrivsFile,
			path:         networkPath,
		},
	}

	var privs []string
	err := json.Unmarshal([]byte(config.VSphereAdminPrivsFile), &privs)
	if err != nil {
		t.Fatalf("failed to validate privs: %v", err)
	}
	vsc.EXPECT().Username().Return("foobar")
	vsc.EXPECT().GetPrivsOnEntity(ctx, networkPath, govmomi.VSphereTypeNetwork, "foobar").Return(privs, nil)

	err = v.validatePrivs(ctx, objects, vsc)
	if err != nil {
		t.Fatalf("failed to validate privs: %v", err)
	}
}

func TestValidatorValidatePrivsError(t *testing.T) {
	v := Validator{}

	ctrl := gomock.NewController(t)
	vsc := mocks.NewMockVSphereClient(ctrl)

	ctx := context.Background()
	networkPath := "/Datacenter/network/path/foo"

	objects := []PrivAssociation{
		{
			objectType:   govmomi.VSphereTypeNetwork,
			privsContent: config.VSphereUserPrivsFile,
			path:         networkPath,
		},
	}

	var privs []string
	err := json.Unmarshal([]byte(config.VSphereAdminPrivsFile), &privs)
	if err != nil {
		t.Fatalf("failed to validate privs: %v", err)
	}
	errMsg := "Could not retrieve privs"
	g := NewWithT(t)
	vsc.EXPECT().Username().Return("foobar")
	vsc.EXPECT().GetPrivsOnEntity(ctx, networkPath, govmomi.VSphereTypeNetwork, "foobar").Return(nil, fmt.Errorf(errMsg))

	err = v.validatePrivs(ctx, objects, vsc)
	g.Expect(err).To(MatchError(ContainSubstring(errMsg)))
}

func TestValidatorValidatePrivsBadJson(t *testing.T) {
	v := Validator{}

	ctrl := gomock.NewController(t)
	vsc := mocks.NewMockVSphereClient(ctrl)
	vsc.EXPECT().Username().Return("foobar")

	ctx := context.Background()
	networkPath := "/Datacenter/network/path/foo"
	g := NewWithT(t)
	errMsg := "invalid character 'h' in literal true (expecting 'r')"

	objects := []PrivAssociation{
		{
			objectType:   govmomi.VSphereTypeNetwork,
			privsContent: "this is bad json",
			path:         networkPath,
		},
	}

	err := v.validatePrivs(ctx, objects, vsc)
	g.Expect(err).To(MatchError(ContainSubstring(errMsg)))
}
