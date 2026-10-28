// SPDX-FileCopyrightText: Contributors to the Gardener project
//
// SPDX-License-Identifier: Apache-2.0

package botanist

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	v1beta1helper "github.com/gardener/gardener/pkg/api/core/v1beta1/helper"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/component"
	extensionsdnsrecord "github.com/gardener/gardener/pkg/component/extensions/dnsrecord"
	"github.com/gardener/gardener/pkg/controllerutils"
	gardenerutils "github.com/gardener/gardener/pkg/utils/gardener"
)

// DefaultExternalDNSRecord creates the default deployer for the external DNSRecord resource.
func (b *Botanist) DefaultExternalDNSRecord() extensionsdnsrecord.Interface {
	values := &extensionsdnsrecord.Values{
		Name:              b.Shoot.GetInfo().Name + "-" + v1beta1constants.DNSRecordExternalName,
		SecretName:        DNSRecordSecretPrefix + "-" + b.Shoot.GetInfo().Name + "-" + v1beta1constants.DNSRecordExternalName,
		Namespace:         b.Shoot.ControlPlaneNamespace,
		TTL:               b.dnsRecordTTLSeconds(),
		AnnotateOperation: controllerutils.HasTask(b.Shoot.GetInfo().Annotations, v1beta1constants.ShootTaskDeployDNSRecordExternal) || b.Shoot.IsRestorePhase(),
		IPStack:           gardenerutils.GetIPStackForShoot(b.Shoot.GetInfo()),
		Labels: map[string]string{
			v1beta1constants.LabelRole:  v1beta1constants.LabelDNSRecordExternal,
			v1beta1constants.GardenRole: v1beta1constants.GardenRoleControlPlane,
		},
	}

	var credentialsDeployer extensionsdnsrecord.CredentialsDeployFunc

	if b.NeedsExternalDNS() {
		values.Type = b.Shoot.ExternalDomain.Provider
		if b.Shoot.ExternalDomain.Zone != "" {
			values.Zone = &b.Shoot.ExternalDomain.Zone
		}
		credentialsDeployer = extensionsdnsrecord.CredentialsDeployerFromCredentials(b.Shoot.ExternalDomain.Credentials, b.Shoot.GetInfo())
		values.DNSName = v1beta1helper.GetAPIServerDomain(*b.Shoot.ExternalClusterDomain)
	}

	return extensionsdnsrecord.New(
		b.Logger,
		b.SeedClientSet.Client(),
		values,
		extensionsdnsrecord.DefaultInterval,
		extensionsdnsrecord.DefaultSevereThreshold,
		extensionsdnsrecord.DefaultTimeout,
		credentialsDeployer,
	)
}

// DefaultPriorExternalDNSRecord creates the default deployer for the DNSRecord resource of the prior external domain.
func (b *Botanist) DefaultPriorExternalDNSRecord() extensionsdnsrecord.Interface {
	values := &extensionsdnsrecord.Values{
		Name:              b.Shoot.GetInfo().Name + "-" + v1beta1constants.DNSRecordPriorExternalName,
		SecretName:        DNSRecordSecretPrefix + "-" + b.Shoot.GetInfo().Name + "-" + v1beta1constants.DNSRecordPriorExternalName,
		Namespace:         b.Shoot.ControlPlaneNamespace,
		TTL:               b.dnsRecordTTLSeconds(),
		AnnotateOperation: controllerutils.HasTask(b.Shoot.GetInfo().Annotations, v1beta1constants.ShootTaskDeployDNSRecordExternal) || b.Shoot.IsRestorePhase(),
		IPStack:           gardenerutils.GetIPStackForShoot(b.Shoot.GetInfo()),
		Labels: map[string]string{
			v1beta1constants.LabelRole:  v1beta1constants.LabelDNSRecordPriorExternal,
			v1beta1constants.GardenRole: v1beta1constants.GardenRoleControlPlane,
		},
	}

	var credentialsDeployer extensionsdnsrecord.CredentialsDeployFunc

	if b.NeedsPriorExternalDNS() {
		values.Type = b.Shoot.PriorExternalDomain.Provider
		if b.Shoot.PriorExternalDomain.Zone != "" {
			values.Zone = &b.Shoot.PriorExternalDomain.Zone
		}
		credentialsDeployer = extensionsdnsrecord.CredentialsDeployerFromCredentials(b.Shoot.PriorExternalDomain.Credentials, b.Shoot.GetInfo())
		values.DNSName = v1beta1helper.GetAPIServerDomain(*b.Shoot.PriorExternalClusterDomain)
	}

	return extensionsdnsrecord.New(
		b.Logger,
		b.SeedClientSet.Client(),
		values,
		extensionsdnsrecord.DefaultInterval,
		extensionsdnsrecord.DefaultSevereThreshold,
		extensionsdnsrecord.DefaultTimeout,
		credentialsDeployer,
	)
}

// DefaultInternalDNSRecord creates the default deployer for the internal DNSRecord resource.
func (b *Botanist) DefaultInternalDNSRecord() extensionsdnsrecord.Interface {
	values := &extensionsdnsrecord.Values{
		Name:                         b.Shoot.GetInfo().Name + "-" + v1beta1constants.DNSRecordInternalName,
		SecretName:                   DNSRecordSecretPrefix + "-" + b.Shoot.GetInfo().Name + "-" + v1beta1constants.DNSRecordInternalName,
		Namespace:                    b.Shoot.ControlPlaneNamespace,
		TTL:                          b.dnsRecordTTLSeconds(),
		ReconcileOnlyOnChangeOrError: b.Shoot.GetInfo().DeletionTimestamp != nil,
		AnnotateOperation: b.Shoot.GetInfo().DeletionTimestamp != nil ||
			controllerutils.HasTask(b.Shoot.GetInfo().Annotations, v1beta1constants.ShootTaskDeployDNSRecordInternal) ||
			b.Shoot.IsRestorePhase(),
		IPStack: gardenerutils.GetIPStackForShoot(b.Shoot.GetInfo()),
		Labels: map[string]string{
			v1beta1constants.LabelRole:  v1beta1constants.LabelDNSRecordInternal,
			v1beta1constants.GardenRole: v1beta1constants.GardenRoleControlPlane,
		},
	}

	var credentialsDeployer extensionsdnsrecord.CredentialsDeployFunc

	if b.NeedsInternalDNS() {
		values.Type = b.Shoot.InternalDomain.Provider
		if b.Shoot.InternalDomain.Zone != "" {
			values.Zone = &b.Shoot.InternalDomain.Zone
		}
		credentialsDeployer = extensionsdnsrecord.CredentialsDeployerFromCredentials(b.Shoot.InternalDomain.Credentials, b.Shoot.GetInfo())
		values.DNSName = v1beta1helper.GetAPIServerDomain(*b.Shoot.InternalClusterDomain)
	}

	return extensionsdnsrecord.New(
		b.Logger,
		b.SeedClientSet.Client(),
		values,
		extensionsdnsrecord.DefaultInterval,
		extensionsdnsrecord.DefaultSevereThreshold,
		extensionsdnsrecord.DefaultTimeout,
		credentialsDeployer,
	)
}

// DefaultPriorInternalDNSRecord creates the default deployer for the DNSRecord resource of the prior internal domain.
func (b *Botanist) DefaultPriorInternalDNSRecord() extensionsdnsrecord.Interface {
	values := &extensionsdnsrecord.Values{
		Name:                         b.Shoot.GetInfo().Name + "-" + v1beta1constants.DNSRecordPriorInternalName,
		SecretName:                   DNSRecordSecretPrefix + "-" + b.Shoot.GetInfo().Name + "-" + v1beta1constants.DNSRecordPriorInternalName,
		Namespace:                    b.Shoot.ControlPlaneNamespace,
		TTL:                          b.dnsRecordTTLSeconds(),
		ReconcileOnlyOnChangeOrError: b.Shoot.GetInfo().DeletionTimestamp != nil,
		AnnotateOperation: b.Shoot.GetInfo().DeletionTimestamp != nil ||
			controllerutils.HasTask(b.Shoot.GetInfo().Annotations, v1beta1constants.ShootTaskDeployDNSRecordInternal) ||
			b.Shoot.IsRestorePhase(),
		IPStack: gardenerutils.GetIPStackForShoot(b.Shoot.GetInfo()),
		Labels: map[string]string{
			v1beta1constants.LabelRole:  v1beta1constants.LabelDNSRecordPriorInternal,
			v1beta1constants.GardenRole: v1beta1constants.GardenRoleControlPlane,
		},
	}

	var credentialsDeployer extensionsdnsrecord.CredentialsDeployFunc

	if b.NeedsPriorInternalDNS() {
		values.Type = b.Shoot.PriorInternalDomain.Provider
		if b.Shoot.PriorInternalDomain.Zone != "" {
			values.Zone = &b.Shoot.PriorInternalDomain.Zone
		}
		credentialsDeployer = extensionsdnsrecord.CredentialsDeployerFromCredentials(b.Shoot.PriorInternalDomain.Credentials, b.Shoot.GetInfo())
		values.DNSName = v1beta1helper.GetAPIServerDomain(*b.Shoot.PriorInternalClusterDomain)
	}

	return extensionsdnsrecord.New(
		b.Logger,
		b.SeedClientSet.Client(),
		values,
		extensionsdnsrecord.DefaultInterval,
		extensionsdnsrecord.DefaultSevereThreshold,
		extensionsdnsrecord.DefaultTimeout,
		credentialsDeployer,
	)
}

// DeployOrDestroyExternalDNSRecord deploys, restores, or destroys the external DNSRecord and waits for the operation to complete.
func (b *Botanist) DeployOrDestroyExternalDNSRecord(ctx context.Context) error {
	if b.NeedsExternalDNS() {
		return b.deployExternalDNSRecord(ctx)
	}
	return b.DestroyExternalDNSRecord(ctx)
}

// DeployOrDestroyPriorExternalDNSRecord deploys, restores, or destroys the prior external DNSRecord and waits for the
// operation to complete.
func (b *Botanist) DeployOrDestroyPriorExternalDNSRecord(ctx context.Context) error {
	if b.NeedsPriorExternalDNS() {
		return b.deployPriorExternalDNSRecord(ctx)
	}
	return b.DestroyPriorExternalDNSRecord(ctx)
}

// DeployOrDestroyPriorInternalDNSRecord deploys, restores, or destroys the prior internal DNSRecord and waits for the
// operation to complete.
func (b *Botanist) DeployOrDestroyPriorInternalDNSRecord(ctx context.Context) error {
	if b.NeedsPriorInternalDNS() {
		if err := b.deployOrRestoreDNSRecord(ctx, b.Shoot.Components.Extensions.PriorInternalDNSRecord); err != nil {
			return err
		}
		return b.Shoot.Components.Extensions.PriorInternalDNSRecord.Wait(ctx)
	}
	return b.DestroyPriorInternalDNSRecord(ctx)
}

// DeployOrDestroyInternalDNSRecord deploys, restores, or destroys the internal DNSRecord and waits for the operation to complete.
func (b *Botanist) DeployOrDestroyInternalDNSRecord(ctx context.Context) error {
	if b.NeedsPriorInternalDNS() {
		if err := b.handOverDNSRecord(ctx, v1beta1constants.DNSRecordInternalName, b.Shoot.InternalClusterDomain, b.MigrateInternalDNSRecord, b.DestroyInternalDNSRecord); err != nil {
			return err
		}
	}

	if b.NeedsInternalDNS() {
		return b.deployInternalDNSRecord(ctx)
	}
	return b.DestroyInternalDNSRecord(ctx)
}

// deployExternalDNSRecord deploys or restores the external DNSRecord and waits for the operation to complete.
func (b *Botanist) deployExternalDNSRecord(ctx context.Context) error {
	if b.NeedsPriorExternalDNS() {
		if err := b.handOverDNSRecord(ctx, v1beta1constants.DNSRecordExternalName, b.Shoot.ExternalClusterDomain, b.MigrateExternalDNSRecord, b.DestroyExternalDNSRecord); err != nil {
			return err
		}
	}

	if err := b.deployOrRestoreDNSRecord(ctx, b.Shoot.Components.Extensions.ExternalDNSRecord); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.ExternalDNSRecord.Wait(ctx)
}

func (b *Botanist) handOverDNSRecord(ctx context.Context, name string, clusterDomain *string, migrate, destroy func(context.Context) error) error {
	dnsRecord := &extensionsv1alpha1.DNSRecord{}
	if err := b.SeedClientSet.Client().Get(ctx, client.ObjectKey{Namespace: b.Shoot.ControlPlaneNamespace, Name: b.Shoot.GetInfo().Name + "-" + name}, dnsRecord); err != nil {
		return client.IgnoreNotFound(err)
	}

	if clusterDomain != nil && dnsRecord.Spec.Name == v1beta1helper.GetAPIServerDomain(*clusterDomain) {
		return nil
	}

	// spec.name of a DNSRecord is immutable. The record of the old domain is migrated, which removes its finalizers but
	// keeps the DNS entry for the prior record. Then it is deleted and recreated for the new domain.
	// The delete flow does not deploy the prior record, thus the record of the old domain is deleted with its DNS entry.
	if b.Shoot.GetInfo().DeletionTimestamp == nil {
		if err := migrate(ctx); err != nil {
			return err
		}
	}
	return destroy(ctx)
}

func (b *Botanist) deployPriorExternalDNSRecord(ctx context.Context) error {
	if err := b.deployOrRestoreDNSRecord(ctx, b.Shoot.Components.Extensions.PriorExternalDNSRecord); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.PriorExternalDNSRecord.Wait(ctx)
}

// deployInternalDNSRecord deploys or restores the internal DNSRecord and waits for the operation to complete.
func (b *Botanist) deployInternalDNSRecord(ctx context.Context) error {
	if err := b.deployOrRestoreDNSRecord(ctx, b.Shoot.Components.Extensions.InternalDNSRecord); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.InternalDNSRecord.Wait(ctx)
}

// DestroyExternalDNSRecord destroys the external DNSRecord and waits for the operation to complete.
func (b *Botanist) DestroyExternalDNSRecord(ctx context.Context) error {
	if err := b.Shoot.Components.Extensions.ExternalDNSRecord.Destroy(ctx); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.ExternalDNSRecord.WaitCleanup(ctx)
}

// DestroyPriorExternalDNSRecord destroys the prior external DNSRecord and waits for the operation to complete.
func (b *Botanist) DestroyPriorExternalDNSRecord(ctx context.Context) error {
	if err := b.Shoot.Components.Extensions.PriorExternalDNSRecord.Destroy(ctx); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.PriorExternalDNSRecord.WaitCleanup(ctx)
}

// DestroyInternalDNSRecord destroys the internal DNSRecord and waits for the operation to complete.
func (b *Botanist) DestroyInternalDNSRecord(ctx context.Context) error {
	if err := b.Shoot.Components.Extensions.InternalDNSRecord.Destroy(ctx); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.InternalDNSRecord.WaitCleanup(ctx)
}

// DestroyPriorInternalDNSRecord destroys the prior internal DNSRecord and waits for the operation to complete.
func (b *Botanist) DestroyPriorInternalDNSRecord(ctx context.Context) error {
	if err := b.Shoot.Components.Extensions.PriorInternalDNSRecord.Destroy(ctx); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.PriorInternalDNSRecord.WaitCleanup(ctx)
}

// MigrateExternalDNSRecord migrates the external DNSRecord and waits for the operation to complete.
func (b *Botanist) MigrateExternalDNSRecord(ctx context.Context) error {
	if err := b.Shoot.Components.Extensions.ExternalDNSRecord.Migrate(ctx); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.ExternalDNSRecord.WaitMigrate(ctx)
}

// MigratePriorExternalDNSRecord migrates the prior external DNSRecord and waits for the operation to complete.
func (b *Botanist) MigratePriorExternalDNSRecord(ctx context.Context) error {
	if err := b.Shoot.Components.Extensions.PriorExternalDNSRecord.Migrate(ctx); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.PriorExternalDNSRecord.WaitMigrate(ctx)
}

// MigrateInternalDNSRecord migrates the internal DNSRecord and waits for the operation to complete.
func (b *Botanist) MigrateInternalDNSRecord(ctx context.Context) error {
	if err := b.Shoot.Components.Extensions.InternalDNSRecord.Migrate(ctx); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.InternalDNSRecord.WaitMigrate(ctx)
}

// MigratePriorInternalDNSRecord migrates the prior internal DNSRecord and waits for the operation to complete.
func (b *Botanist) MigratePriorInternalDNSRecord(ctx context.Context) error {
	if err := b.Shoot.Components.Extensions.PriorInternalDNSRecord.Migrate(ctx); err != nil {
		return err
	}
	return b.Shoot.Components.Extensions.PriorInternalDNSRecord.WaitMigrate(ctx)
}

func (b *Botanist) deployOrRestoreDNSRecord(ctx context.Context, dnsRecord component.DeployMigrateWaiter) error {
	if b.Shoot.IsRestorePhase() {
		return dnsRecord.Restore(ctx, b.Shoot.GetShootState())
	}
	return dnsRecord.Deploy(ctx)
}

func (b *Botanist) dnsRecordTTLSeconds() *int64 {
	if b.Config != nil && b.Config.Controllers != nil && b.Config.Controllers.Shoot != nil {
		return b.Config.Controllers.Shoot.DNSEntryTTLSeconds
	}
	return new(int64(120))
}
