// +build !ignore_autogenerated_openshift

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package api

import (
	pkg_api "k8s.io/kubernetes/pkg/api"
	unversioned "k8s.io/kubernetes/pkg/api/unversioned"
	conversion "k8s.io/kubernetes/pkg/conversion"
	runtime "k8s.io/kubernetes/pkg/runtime"
	reflect "reflect"
)

func init() {
	SchemeBuilder.Register(RegisterDeepCopies)
}

// RegisterDeepCopies adds deep-copy functions to the given scheme. Public
// to allow building arbitrary schemes.
func RegisterDeepCopies(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedDeepCopyFuncs(
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_AppliedClusterResourceQuota, InType: reflect.TypeOf(&AppliedClusterResourceQuota{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_AppliedClusterResourceQuotaList, InType: reflect.TypeOf(&AppliedClusterResourceQuotaList{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_ClusterResourceQuota, InType: reflect.TypeOf(&ClusterResourceQuota{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_ClusterResourceQuotaList, InType: reflect.TypeOf(&ClusterResourceQuotaList{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_ClusterResourceQuotaSelector, InType: reflect.TypeOf(&ClusterResourceQuotaSelector{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_ClusterResourceQuotaSpec, InType: reflect.TypeOf(&ClusterResourceQuotaSpec{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_ClusterResourceQuotaStatus, InType: reflect.TypeOf(&ClusterResourceQuotaStatus{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_api_ResourceQuotasStatusByNamespace, InType: reflect.TypeOf(&ResourceQuotasStatusByNamespace{})},
	)
}

func DeepCopy_api_AppliedClusterResourceQuota(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*AppliedClusterResourceQuota)
		out := out.(*AppliedClusterResourceQuota)
		out.TypeMeta = in.TypeMeta
		if err := pkg_api.DeepCopy_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, c); err != nil {
			return err
		}
		if err := DeepCopy_api_ClusterResourceQuotaSpec(&in.Spec, &out.Spec, c); err != nil {
			return err
		}
		if err := DeepCopy_api_ClusterResourceQuotaStatus(&in.Status, &out.Status, c); err != nil {
			return err
		}
		return nil
	}
}

func DeepCopy_api_AppliedClusterResourceQuotaList(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*AppliedClusterResourceQuotaList)
		out := out.(*AppliedClusterResourceQuotaList)
		out.TypeMeta = in.TypeMeta
		out.ListMeta = in.ListMeta
		if in.Items != nil {
			in, out := &in.Items, &out.Items
			*out = make([]AppliedClusterResourceQuota, len(*in))
			for i := range *in {
				if err := DeepCopy_api_AppliedClusterResourceQuota(&(*in)[i], &(*out)[i], c); err != nil {
					return err
				}
			}
		} else {
			out.Items = nil
		}
		return nil
	}
}

func DeepCopy_api_ClusterResourceQuota(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*ClusterResourceQuota)
		out := out.(*ClusterResourceQuota)
		out.TypeMeta = in.TypeMeta
		if err := pkg_api.DeepCopy_api_ObjectMeta(&in.ObjectMeta, &out.ObjectMeta, c); err != nil {
			return err
		}
		if err := DeepCopy_api_ClusterResourceQuotaSpec(&in.Spec, &out.Spec, c); err != nil {
			return err
		}
		if err := DeepCopy_api_ClusterResourceQuotaStatus(&in.Status, &out.Status, c); err != nil {
			return err
		}
		return nil
	}
}

func DeepCopy_api_ClusterResourceQuotaList(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*ClusterResourceQuotaList)
		out := out.(*ClusterResourceQuotaList)
		out.TypeMeta = in.TypeMeta
		out.ListMeta = in.ListMeta
		if in.Items != nil {
			in, out := &in.Items, &out.Items
			*out = make([]ClusterResourceQuota, len(*in))
			for i := range *in {
				if err := DeepCopy_api_ClusterResourceQuota(&(*in)[i], &(*out)[i], c); err != nil {
					return err
				}
			}
		} else {
			out.Items = nil
		}
		return nil
	}
}

func DeepCopy_api_ClusterResourceQuotaSelector(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*ClusterResourceQuotaSelector)
		out := out.(*ClusterResourceQuotaSelector)
		if in.LabelSelector != nil {
			in, out := &in.LabelSelector, &out.LabelSelector
			*out = new(unversioned.LabelSelector)
			if err := unversioned.DeepCopy_unversioned_LabelSelector(*in, *out, c); err != nil {
				return err
			}
		} else {
			out.LabelSelector = nil
		}
		if in.AnnotationSelector != nil {
			in, out := &in.AnnotationSelector, &out.AnnotationSelector
			*out = make(map[string]string)
			for key, val := range *in {
				(*out)[key] = val
			}
		} else {
			out.AnnotationSelector = nil
		}
		return nil
	}
}

func DeepCopy_api_ClusterResourceQuotaSpec(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*ClusterResourceQuotaSpec)
		out := out.(*ClusterResourceQuotaSpec)
		if err := DeepCopy_api_ClusterResourceQuotaSelector(&in.Selector, &out.Selector, c); err != nil {
			return err
		}
		if err := pkg_api.DeepCopy_api_ResourceQuotaSpec(&in.Quota, &out.Quota, c); err != nil {
			return err
		}
		return nil
	}
}

func DeepCopy_api_ClusterResourceQuotaStatus(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*ClusterResourceQuotaStatus)
		out := out.(*ClusterResourceQuotaStatus)
		if err := pkg_api.DeepCopy_api_ResourceQuotaStatus(&in.Total, &out.Total, c); err != nil {
			return err
		}
		if err := DeepCopy_api_ResourceQuotasStatusByNamespace(&in.Namespaces, &out.Namespaces, c); err != nil {
			return err
		}
		return nil
	}
}

func DeepCopy_api_ResourceQuotasStatusByNamespace(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*ResourceQuotasStatusByNamespace)
		out := out.(*ResourceQuotasStatusByNamespace)
		if newVal, err := c.DeepCopy(&in.orderedMap); err != nil {
			return err
		} else {
			out.orderedMap = *newVal.(*orderedMap)
		}
		return nil
	}
}
