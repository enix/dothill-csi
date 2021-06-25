/*
 * Copyright (c) 2021 Enix, SAS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 *
 * Authors:
 * Paul Laffitte <paul.laffitte@enix.fr>
 * Arthur Chaloin <arthur.chaloin@enix.fr>
 * Alexandre Buisine <alexandre.buisine@enix.fr>
 */

package controller

import (
	"context"
	"fmt"
	"strconv"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/enix/san-iscsi-csi/pkg/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ControllerExpandVolume expands a volume to the given new size
func (controller *Controller) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	volumeID := req.GetVolumeId()
	if volumeID == "" {
		return nil, status.Error(codes.InvalidArgument, "cannot expand a volume with an empty ID")
	}

	common.AddLogTag(ctx, "volumeId", volumeID)

	newSize := req.GetCapacityRange().GetRequiredBytes()
	if newSize == 0 {
		newSize = req.GetCapacityRange().GetLimitBytes()
	}

	response, _, err := controller.dothillClient.ShowVolumes(volumeID)
	var expansionSize int64
	if err != nil {
		return nil, err
	} else if volume, ok := response.ObjectsMap["volume"]; !ok {
		return nil, fmt.Errorf("volume %q not found", volumeID)
	} else if sizeNumeric, ok := volume.PropertiesMap["size-numeric"]; !ok {
		return nil, fmt.Errorf("could not get current volume size, thus volume expansion is not possible")
	} else if currentBlocks, err := strconv.ParseInt(sizeNumeric.Data, 10, 32); err != nil {
		return nil, fmt.Errorf("could not parse volume size: %v", err)
	} else {
		currentSize := currentBlocks * 512
		expansionSize = newSize - currentSize
		common.LogInfoS(ctx, "expanding volume", "oldSize", currentSize, "size", newSize, "addedBytes", expansionSize)
	}

	expansionSizeStr := getSizeStr(expansionSize)
	if _, _, err := controller.dothillClient.ExpandVolume(volumeID, expansionSizeStr); err != nil {
		return nil, err
	}

	common.LogInfoS(ctx, "volume successfully expanded")

	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         newSize,
		NodeExpansionRequired: true,
	}, nil
}
