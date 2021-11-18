package cli

import (
	"fmt"
)

const AWS_TAG_KEY = "lakitu";

func createVolumeAndAttach(shell *Shell, snapshotId string) error {
	fmt.Println("Creating volume from snapshot and attaching to instance...")

	_, _, err := shell.execute(fmt.Sprintf(`
		Function Get-DeviceName() {
		    param(
				[parameter (mandatory=$true)][string] $InstanceId            
			)            
		    $CurrentDevice = Get-EC2InstanceAttribute $InstanceId -Attribute blockDeviceMapping | Select-Object -ExpandProperty BlockDeviceMappings | Select-Object -last 1            
		If ($CurrentDevice.DeviceName -eq '/dev/sda1') {            
		    $NewDevice = 'xvdf'            
		    return $NewDevice            
		}            
		Else {            
		    $a = $CurrentDevice.DeviceName.ToCharArray()            
		    $inc = +1            
		    $a[3] = [char]([int]($a[3])+$inc)            
		    $NewDevice = -join $a            
		    return $NewDevice            
		    }
		}

		$instanceId = Get-EC2InstanceMetadata -Path "/instance-id";
		$zone = Get-EC2InstanceMetadata -Path "/placement/availability-zone";
		$volumeType = "gp3";

		$TagSpecification = [Amazon.EC2.Model.TagSpecification]::new()
		$TagSpecification.ResourceType = 'Volume'
    	$TagSpecification.Tags.Add([Amazon.EC2.Model.Tag]@{
	        Key   = "%s"
	        Value = ""
	    });

		$volumeId = (New-EC2Volume -SnapshotId "%s" -AvailabilityZone $zone -VolumeType $volumeType -TagSpecification $TagSpecification).VolumeId;

		while ((Get-EC2Volume -VolumeId $volumeId).State -ne "available") {
		  write-host -NoNewline "."
		  sleep 5
		}

		$deviceName = Get-DeviceName $InstanceId;
		Add-EC2Volume -InstanceId $instanceId -VolumeId $volumeId -Device $deviceName;

		while ((Get-EC2Volume -VolumeId $volumeId).State -ne "in-use") {
		  write-host -NoNewline "."
		  sleep 5
		}

		$diskNum = (Get-Disk | Where OperationalStatus -eq 'Offline').Number;
		Set-Disk -Number $diskNum -IsOffline $False;
	`, AWS_TAG_KEY, snapshotId))

	if err != nil {
		return err
	}

	return nil
}

func setupInstanceStore(shell *Shell) error {
	fmt.Println("Mounting instance store if exists...")

	_, _, err := shell.execute(`
		Get-Disk |
		Where partitionstyle -eq 'raw'  | 
		Where SerialNumber -like 'AWS*' | 
		Initialize-Disk -PartitionStyle MBR -PassThru | 
		New-Partition -DriveLetter D -UseMaximumSize  | 
		Format-Volume -FileSystem NTFS -NewFileSystemLabel "InstanceStore" -Confirm:$false
	`)

	if err != nil {
		return err
	}

	return nil
}

func setupNewVolume(shell *Shell) error {
	fmt.Println("Mounting new volume if exists...")

	_, _, err := shell.execute(`
		Get-Disk |
		Where partitionstyle -eq 'raw'  | 
		Where SerialNumber -like 'vol*' | 
		Initialize-Disk -PartitionStyle MBR -PassThru | 
		New-Partition -DriveLetter L -UseMaximumSize  | 
		Format-Volume -FileSystem NTFS -NewFileSystemLabel "LakituStorage" -Confirm:$false
	`)

	if err != nil {
		return err
	}

	return nil
}

func MountSnapshot(param string) error {
	ps := NewShell()

	err := setupInstanceStore(ps)
	if err != nil {
		return err
	}

	if param == "" {
		return nil
	}

	if param == "new" {
		return setupNewVolume(ps)
	} else {
		return createVolumeAndAttach(ps, param)
	}

	return nil
}
