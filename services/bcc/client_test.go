package bcc

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/baidubce/bce-sdk-go/model"
	"github.com/baidubce/bce-sdk-go/services/bcc/api"
	"github.com/baidubce/bce-sdk-go/util/log"
)

var (
	BCC_CLIENT              *Client
	BCC_TestCdsId           string
	BCC_TestBccId           string
	BCC_TestSecurityGroupId string
	BCC_TestImageId         string
	BCC_TestSnapshotId      string
	BCC_TestAspId           string
	BCC_TestDeploySetId     string
)

// For security reason, ak/sk should not hard write here.
type Conf struct {
	AK       string
	SK       string
	Endpoint string
}

func init() {
	_, f, _, _ := runtime.Caller(0)
	for i := 0; i < 7; i++ {
		f = filepath.Dir(f)
	}
	conf := filepath.Join(f, "/go_conf.json")
	fmt.Println(conf)
	fp, err := os.Open(conf)
	if err != nil {
		log.Fatal("config json file of ak/sk not given:", conf)
		os.Exit(1)
	}
	decoder := json.NewDecoder(fp)
	confObj := &Conf{}
	decoder.Decode(confObj)

	BCC_CLIENT, _ = NewClient(confObj.AK, confObj.SK, confObj.Endpoint)
	log.SetLogLevel(log.WARN)
	//log.SetLogLevel(log.DEBUG)
}

// ExpectEqual is the helper function for test each case
func ExpectEqual(alert func(format string, args ...interface{}),
	expected interface{}, actual interface{}) bool {
	expectedValue, actualValue := reflect.ValueOf(expected), reflect.ValueOf(actual)
	equal := false
	switch {
	case expected == nil && actual == nil:
		return true
	case expected != nil && actual == nil:
		equal = expectedValue.IsNil()
	case expected == nil && actual != nil:
		equal = actualValue.IsNil()
	default:
		if actualType := reflect.TypeOf(actual); actualType != nil {
			if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
				equal = reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual)
			}
		}
	}
	if !equal {
		_, file, line, _ := runtime.Caller(1)
		alert("%s:%d: missmatch, expect %v but %v", file, line, expected, actual)
		return false
	}
	return true
}

func TestCreateInstance(t *testing.T) {
	InternalIps := []string{"ip"}
	createInstanceArgs := &api.CreateInstanceArgs{
		ImageId: "ImageId",
		Billing: api.Billing{
			PaymentTiming: api.PaymentTimingPostPaid,
		},
		InstanceType:        api.InstanceTypeN3,
		CpuCount:            1,
		MemoryCapacityInGB:  4,
		RootDiskSizeInGb:    40,
		RootDiskStorageType: api.StorageTypeHP1,
		ZoneName:            "ZoneName",
		SubnetId:            "SubnetId",
		SecurityGroupId:     "SecurityGroupId",
		RelationTag:         true,
		PurchaseCount:       1,
		Name:                "sdkTest",
		KeypairId:           "KeypairId",
		InternalIps:         InternalIps,
	}
	createResult, err := BCC_CLIENT.CreateInstance(createInstanceArgs)
	ExpectEqual(t.Errorf, err, nil)
	BCC_TestBccId = createResult.InstanceIds[0]
}

func TestCreateInstanceBySpec(t *testing.T) {
	createInstanceBySpecArgs := &api.CreateInstanceBySpecArgs{
		ImageId:   "m-1PyVLtic",
		Spec:      "bcc.g2.c2m8",
		Name:      "sdkTest2",
		AdminPass: "123qaz!@#",
		ZoneName:  "cn-bj-a",
		Billing: api.Billing{
			PaymentTiming: api.PaymentTimingPostPaid,
		},
	}
	_, err := BCC_CLIENT.CreateInstanceBySpec(createInstanceBySpecArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestCreateSecurityGroup(t *testing.T) {
	args := &api.CreateSecurityGroupArgs{
		Name: "testSecurityGroup",
		Rules: []api.SecurityGroupRuleModel{
			{
				Remark:        "备注",
				Protocol:      "tcp",
				PortRange:     "1-65535",
				Direction:     "ingress",
				SourceIp:      "",
				SourceGroupId: "",
			},
		},
	}
	result, err := BCC_CLIENT.CreateSecurityGroup(args)
	ExpectEqual(t.Errorf, err, nil)
	BCC_TestSecurityGroupId = result.SecurityGroupId
}

func TestCreateImage(t *testing.T) {
	args := &api.CreateImageArgs{
		ImageName:  "testImageName",
		InstanceId: BCC_TestBccId,
	}
	result, err := BCC_CLIENT.CreateImage(args)
	ExpectEqual(t.Errorf, err, nil)
	BCC_TestImageId = result.ImageId
}

func TestListInstances(t *testing.T) {
	listArgs := &api.ListInstanceArgs{}
	_, err := BCC_CLIENT.ListInstances(listArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestGetInstanceDetail(t *testing.T) {
	_, err := BCC_CLIENT.GetInstanceDetail(BCC_TestBccId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestResizeInstance(t *testing.T) {
	resizeArgs := &api.ResizeInstanceArgs{
		CpuCount:           2,
		MemoryCapacityInGB: 4,
	}
	err := BCC_CLIENT.ResizeInstance(BCC_TestBccId, resizeArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestStopInstanceWithNoCharge(t *testing.T) {
	err := BCC_CLIENT.StopInstanceWithNoCharge(BCC_TestBccId, true, true)
	ExpectEqual(t.Errorf, err, nil)
}

func TestStopInstance(t *testing.T) {
	err := BCC_CLIENT.StopInstance(BCC_TestBccId, true)
	ExpectEqual(t.Errorf, err, nil)
}

func TestStartInstance(t *testing.T) {
	err := BCC_CLIENT.StartInstance(BCC_TestBccId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestRebootInstance(t *testing.T) {
	err := BCC_CLIENT.RebootInstance(BCC_TestBccId, true)
	ExpectEqual(t.Errorf, err, nil)
}

func TestRebuildInstance(t *testing.T) {
	rebuildArgs := &api.RebuildInstanceArgs{
		ImageId:   "m-DpgNg8lO",
		AdminPass: "123qaz!@#",
	}
	err := BCC_CLIENT.RebuildInstance(BCC_TestBccId, rebuildArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestChangeInstancePass(t *testing.T) {
	changeArgs := &api.ChangeInstancePassArgs{
		AdminPass: "321zaq#@!",
	}
	err := BCC_CLIENT.ChangeInstancePass(BCC_TestBccId, changeArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestModifyInstanceAttribute(t *testing.T) {
	modifyArgs := &api.ModifyInstanceAttributeArgs{
		Name: "test-modify",
	}
	err := BCC_CLIENT.ModifyInstanceAttribute(BCC_TestBccId, modifyArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestModifyInstanceDesc(t *testing.T) {
	modifyArgs := &api.ModifyInstanceDescArgs{
		Description: "test modify",
	}
	err := BCC_CLIENT.ModifyInstanceDesc(BCC_TestBccId, modifyArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestGetInstanceVNC(t *testing.T) {
	_, err := BCC_CLIENT.GetInstanceVNC(BCC_TestBccId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestBatchAddIp(t *testing.T) {
	privateIps := []string{"192.168.16.17"}
	batchAddIpArgs := &api.BatchAddIpArgs{
		InstanceId: BCC_TestBccId,
		PrivateIps: privateIps,
	}
	err := BCC_CLIENT.BatchAddIP(batchAddIpArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestBatchDelIp(t *testing.T) {
	privateIps := []string{"192.168.16.17"}
	batchDelIpArgs := &api.BatchDelIpArgs{
		InstanceId: BCC_TestBccId,
		PrivateIps: privateIps,
	}
	err := BCC_CLIENT.BatchDelIP(batchDelIpArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestBindSecurityGroup(t *testing.T) {
	err := BCC_CLIENT.BindSecurityGroup(BCC_TestBccId, BCC_TestSecurityGroupId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestUnBindSecurityGroup(t *testing.T) {
	err := BCC_CLIENT.UnBindSecurityGroup(BCC_TestBccId, BCC_TestSecurityGroupId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestCreateCDSVolume(t *testing.T) {
	args := &api.CreateCDSVolumeArgs{
		PurchaseCount: 1,
		CdsSizeInGB:   40,
		StorageType:   api.StorageTypeSSD,
		Billing: &api.Billing{
			PaymentTiming: api.PaymentTimingPostPaid,
		},
		EncryptKey: "EncryptKey",
	}

	result, err := BCC_CLIENT.CreateCDSVolume(args)
	ExpectEqual(t.Errorf, err, nil)
	BCC_TestCdsId = result.VolumeIds[0]
}

func TestCreateSnapshot(t *testing.T) {
	args := &api.CreateSnapshotArgs{
		VolumeId:     BCC_TestCdsId,
		SnapshotName: "testSnapshotName",
	}
	result, err := BCC_CLIENT.CreateSnapshot(args)
	ExpectEqual(t.Errorf, err, nil)

	BCC_TestSnapshotId = result.SnapshotId
}

func TestListSnapshot(t *testing.T) {
	args := &api.ListSnapshotArgs{}
	_, err := BCC_CLIENT.ListSnapshot(args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestGetSnapshotDetail(t *testing.T) {
	_, err := BCC_CLIENT.GetSnapshotDetail(BCC_TestSnapshotId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestDeleteSnapshot(t *testing.T) {
	err := BCC_CLIENT.DeleteSnapshot(BCC_TestSnapshotId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestCreateAutoSnapshotPolicy(t *testing.T) {
	args := &api.CreateASPArgs{
		Name:           "testAspName",
		TimePoints:     []string{"20"},
		RepeatWeekdays: []string{"1", "5"},
		RetentionDays:  "7",
	}
	result, err := BCC_CLIENT.CreateAutoSnapshotPolicy(args)
	ExpectEqual(t.Errorf, err, nil)
	BCC_TestAspId = result.AspId
}

func TestAttachAutoSnapshotPolicy(t *testing.T) {
	args := &api.AttachASPArgs{
		VolumeIds: []string{BCC_TestCdsId},
	}
	err := BCC_CLIENT.AttachAutoSnapshotPolicy(BCC_TestAspId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestDetachAutoSnapshotPolicy(t *testing.T) {
	args := &api.DetachASPArgs{
		VolumeIds: []string{BCC_TestCdsId},
	}
	err := BCC_CLIENT.DetachAutoSnapshotPolicy(BCC_TestAspId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestListAutoSnapshotPolicy(t *testing.T) {
	args := &api.ListASPArgs{}
	_, err := BCC_CLIENT.ListAutoSnapshotPolicy(args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestGetAutoSnapshotPolicy(t *testing.T) {
	_, err := BCC_CLIENT.GetAutoSnapshotPolicy(BCC_TestAspId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestDeleteAutoSnapshotPolicy(t *testing.T) {
	err := BCC_CLIENT.DeleteAutoSnapshotPolicy(BCC_TestAspId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestListCDSVolume(t *testing.T) {
	queryArgs := &api.ListCDSVolumeArgs{}
	_, err := BCC_CLIENT.ListCDSVolume(queryArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestGetCDSVolumeDetail(t *testing.T) {
	_, err := BCC_CLIENT.GetCDSVolumeDetail(BCC_TestCdsId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestAttachCDSVolume(t *testing.T) {
	args := &api.AttachVolumeArgs{
		InstanceId: BCC_TestBccId,
	}

	_, err := BCC_CLIENT.AttachCDSVolume(BCC_TestCdsId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestDetachCDSVolume(t *testing.T) {
	args := &api.DetachVolumeArgs{
		InstanceId: BCC_TestBccId,
	}

	err := BCC_CLIENT.DetachCDSVolume(BCC_TestCdsId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestResizeCDSVolume(t *testing.T) {
	args := &api.ResizeCSDVolumeArgs{
		NewCdsSizeInGB: 100,
		NewVolumeType:  api.StorageTypeHdd,
	}

	err := BCC_CLIENT.ResizeCDSVolume(BCC_TestCdsId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestPurchaseReservedCDSVolume(t *testing.T) {
	args := &api.PurchaseReservedCSDVolumeArgs{
		Billing: &api.Billing{
			Reservation: &api.Reservation{
				ReservationLength:   1,
				ReservationTimeUnit: "Month",
			},
		},
	}

	err := BCC_CLIENT.PurchaseReservedCDSVolume(BCC_TestCdsId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestRenameCDSVolume(t *testing.T) {
	args := &api.RenameCSDVolumeArgs{
		Name: "testRenamedName",
	}

	err := BCC_CLIENT.RenameCDSVolume(BCC_TestCdsId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestModifyCDSVolume(t *testing.T) {
	args := &api.ModifyCSDVolumeArgs{
		CdsName: "modifiedName",
		Desc:    "modifiedDesc",
	}

	err := BCC_CLIENT.ModifyCDSVolume(BCC_TestCdsId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestModifyChargeTypeCDSVolume(t *testing.T) {
	args := &api.ModifyChargeTypeCSDVolumeArgs{
		Billing: &api.Billing{
			Reservation: &api.Reservation{
				ReservationLength: 1,
			},
		},
	}

	err := BCC_CLIENT.ModifyChargeTypeCDSVolume(BCC_TestCdsId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestDeleteCDSVolumeNew(t *testing.T) {
	args := &api.DeleteCDSVolumeArgs{
		AutoSnapshot: "on",
	}

	err := BCC_CLIENT.DeleteCDSVolumeNew(BCC_TestCdsId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestDeleteCDSVolume(t *testing.T) {
	err := BCC_CLIENT.DeleteCDSVolume(BCC_TestCdsId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestListSecurityGroup(t *testing.T) {
	queryArgs := &api.ListSecurityGroupArgs{}
	_, err := BCC_CLIENT.ListSecurityGroup(queryArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestAuthorizeSecurityGroupRule(t *testing.T) {
	args := &api.AuthorizeSecurityGroupArgs{
		Rule: &api.SecurityGroupRuleModel{
			Remark:        "备注",
			Protocol:      "udp",
			PortRange:     "1-65535",
			Direction:     "ingress",
			SourceIp:      "",
			SourceGroupId: "",
		},
	}
	err := BCC_CLIENT.AuthorizeSecurityGroupRule(BCC_TestSecurityGroupId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestRevokeSecurityGroupRule(t *testing.T) {
	args := &api.RevokeSecurityGroupArgs{
		Rule: &api.SecurityGroupRuleModel{
			Remark:        "备注",
			Protocol:      "udp",
			PortRange:     "1-65535",
			Direction:     "ingress",
			SourceIp:      "",
			SourceGroupId: "",
		},
	}
	err := BCC_CLIENT.RevokeSecurityGroupRule(BCC_TestSecurityGroupId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestDeleteSecurityGroupRule(t *testing.T) {
	err := BCC_CLIENT.DeleteSecurityGroup(BCC_TestSecurityGroupId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestListImage(t *testing.T) {
	queryArgs := &api.ListImageArgs{}
	_, err := BCC_CLIENT.ListImage(queryArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestGetImageDetail(t *testing.T) {
	_, err := BCC_CLIENT.GetImageDetail(BCC_TestImageId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestRemoteCopyImage(t *testing.T) {
	args := &api.RemoteCopyImageArgs{
		Name:       "testRemoteCopy",
		DestRegion: []string{"bj"},
	}
	err := BCC_CLIENT.RemoteCopyImage(BCC_TestImageId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestCancelRemoteCopyImage(t *testing.T) {
	err := BCC_CLIENT.CancelRemoteCopyImage(BCC_TestImageId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestShareImage(t *testing.T) {
	args := &api.SharedUser{
		AccountId: "id",
	}
	err := BCC_CLIENT.ShareImage(BCC_TestImageId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestUnShareImage(t *testing.T) {
	args := &api.SharedUser{
		AccountId: "id",
	}
	err := BCC_CLIENT.UnShareImage(BCC_TestImageId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestGetImageSharedUser(t *testing.T) {
	_, err := BCC_CLIENT.GetImageSharedUser(BCC_TestImageId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestGetImageOS(t *testing.T) {
	args := &api.GetImageOsArgs{}
	_, err := BCC_CLIENT.GetImageOS(args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestDeleteImage(t *testing.T) {
	err := BCC_CLIENT.DeleteImage(BCC_TestImageId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestDeleteInstance(t *testing.T) {
	err := BCC_CLIENT.DeleteInstance(BCC_TestBccId)
	ExpectEqual(t.Errorf, err, nil)
}

func TestListSpec(t *testing.T) {
	_, err := BCC_CLIENT.ListSpec()
	ExpectEqual(t.Errorf, err, nil)
}

func TestListZone(t *testing.T) {
	_, err := BCC_CLIENT.ListZone()
	ExpectEqual(t.Errorf, err, nil)
}

func TestCreateDeploySet(t *testing.T) {
	testDeploySetName := "testName"
	testDeployDesc := "testDesc"
	testStrategy := "HOST_HA"
	queryArgs := &api.CreateDeploySetArgs{
		Strategy: testStrategy,
		Name:     testDeploySetName,
		Desc:     testDeployDesc,
	}
	rep, err := BCC_CLIENT.CreateDeploySet(queryArgs)
	fmt.Println(rep)
	ExpectEqual(t.Errorf, err, nil)
	BCC_TestDeploySetId = rep.DeploySetIds[0]
}

func TestListDeploySets(t *testing.T) {
	rep, err := BCC_CLIENT.ListDeploySets()
	fmt.Println(rep)
	ExpectEqual(t.Errorf, err, nil)
}

func TestModifyDeploySet(t *testing.T) {
	testDeploySetName := "testName"
	testDeployDesc := "goDesc"
	queryArgs := &api.ModifyDeploySetArgs{
		Name: testDeploySetName,
		Desc: testDeployDesc,
	}
	BCC_TestDeploySetId = "DeploySetId"
	rep, err := BCC_CLIENT.ModifyDeploySet(BCC_TestDeploySetId, queryArgs)
	fmt.Println(rep)
	ExpectEqual(t.Errorf, err, nil)
}

func TestDeleteDeploySet(t *testing.T) {
	testDeleteDeploySetId := "DeploySetId"
	err := BCC_CLIENT.DeleteDeploySet(testDeleteDeploySetId)
	fmt.Println(err)
	ExpectEqual(t.Errorf, err, nil)
}

func TestResizeInstanceBySpec(t *testing.T) {
	resizeArgs := &api.ResizeInstanceArgs{
		Spec: "Spec",
	}
	err := BCC_CLIENT.ResizeInstanceBySpec(BCC_TestBccId, resizeArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestBatchRebuildInstances(t *testing.T) {
	rebuildArgs := &api.RebuildBatchInstanceArgs{
		ImageId:     "ImageId",
		AdminPass:   "123qaz!@#",
		InstanceIds: []string{BCC_TestBccId},
	}
	err := BCC_CLIENT.BatchRebuildInstances(rebuildArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestChangeToPrepaid(t *testing.T) {
	args := &api.ChangeToPrepaidRequest{
		Duration:    1,
		RelationCds: true,
	}
	_, err := BCC_CLIENT.ChangeToPrepaid(BCC_TestBccId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestBindInstanceToTags(t *testing.T) {
	args := &api.BindTagsRequest{
		ChangeTags: []model.TagModel{
			{
				TagKey:   "TagKey",
				TagValue: "TagValue",
			},
		},
	}
	err := BCC_CLIENT.BindInstanceToTags(BCC_TestBccId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestUnBindInstanceToTags(t *testing.T) {
	args := &api.UnBindTagsRequest{
		ChangeTags: []model.TagModel{
			{
				TagKey:   "TagKey",
				TagValue: "TagValue",
			},
		},
	}
	err := BCC_CLIENT.UnBindInstanceToTags(BCC_TestBccId, args)
	ExpectEqual(t.Errorf, err, nil)
}

func TestGetInstanceNoChargeList(t *testing.T) {
	listArgs := &api.ListInstanceArgs{}
	_, err := BCC_CLIENT.GetInstanceNoChargeList(listArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestCreateBidInstance(t *testing.T) {
	createInstanceArgs := &api.CreateInstanceArgs{
		ImageId: "ImageId",
		Billing: api.Billing{
			PaymentTiming: api.PaymentTimingBidding,
		},
		InstanceType:        api.InstanceTypeN3,
		CpuCount:            1,
		MemoryCapacityInGB:  4,
		RootDiskSizeInGb:    40,
		RootDiskStorageType: api.StorageTypeHP1,
		ZoneName:            "zoneName",
		SubnetId:            "SubnetId",
		SecurityGroupId:     "SecurityGroupId",
		RelationTag:         true,
		PurchaseCount:       1,
		Name:                "sdkTest",
		BidModel:            "BidModel",
		BidPrice:            "BidPrice",
	}
	createResult, err := BCC_CLIENT.CreateBidInstance(createInstanceArgs)
	ExpectEqual(t.Errorf, err, nil)
	BCC_TestBccId = createResult.InstanceIds[0]
}

func TestCancelBidOrder(t *testing.T) {
	createInstanceArgs := &api.CancelBidOrderRequest{
		OrderId: "OrderId",
	}
	_, err := BCC_CLIENT.CancelBidOrder(createInstanceArgs)
	ExpectEqual(t.Errorf, err, nil)
}

func TestInstancePurchaseReserved(t *testing.T) {
	purchaseReservedArgs := &api.PurchaseReservedArgs{
		Billing: api.Billing{
			PaymentTiming: api.PaymentTimingPrePaid,
			Reservation: &api.Reservation{
				ReservationLength: 1,
			},
		},
		RelatedRenewFlag: "CDS",
	}
	err := BCC_CLIENT.InstancePurchaseReserved(BCC_TestBccId, purchaseReservedArgs)
	ExpectEqual(t.Errorf, err, nil)
}
