package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	autoscalingv1a1 "github.com/kubewharf/katalyst-api/pkg/apis/autoscaling/v1alpha1"
	autoscalingv1a2 "github.com/kubewharf/katalyst-api/pkg/apis/autoscaling/v1alpha2"
	configv1a1 "github.com/kubewharf/katalyst-api/pkg/apis/config/v1alpha1"
	nodev1a1 "github.com/kubewharf/katalyst-api/pkg/apis/node/v1alpha1"
	overcommitv1a1 "github.com/kubewharf/katalyst-api/pkg/apis/overcommit/v1alpha1"
	recommendationv1a1 "github.com/kubewharf/katalyst-api/pkg/apis/recommendation/v1alpha1"
	schedcfgv1b3 "github.com/kubewharf/katalyst-api/pkg/apis/scheduling/config/v1beta3"
	tidev1a1 "github.com/kubewharf/katalyst-api/pkg/apis/tide/v1alpha1"
	workloadv1a1 "github.com/kubewharf/katalyst-api/pkg/apis/workload/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
)

func main() {
	scheme := runtime.NewScheme()
	configv1a1.AddToScheme(scheme)
	autoscalingv1a1.AddToScheme(scheme)
	autoscalingv1a2.AddToScheme(scheme)
	nodev1a1.AddToScheme(scheme)
	overcommitv1a1.AddToScheme(scheme)
	recommendationv1a1.AddToScheme(scheme)
	schedcfgv1b3.AddToScheme(scheme)
	tidev1a1.AddToScheme(scheme)
	workloadv1a1.AddToScheme(scheme)

	seenTypes := make(map[reflect.Type]struct{})
	var errs []error
	for _, typ := range scheme.AllKnownTypes() {
		fmt.Printf("Checking %s.%s\n", typ.PkgPath(), typ.Name())
		checkType(typ, typ.Name(), seenTypes, &errs)
	}

	for _, err := range errs {
		fmt.Println(err)
	}

	if len(errs) > 0 {
		os.Exit(1)
	}
}

func parseTag(tag string) (name string) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	} else {
		return tag
	}
}

// typ must be a struct type
func checkType(typ reflect.Type, path string, seenTypes map[reflect.Type]struct{}, errs *[]error) {
	if _, ok := seenTypes[typ]; ok {
		return
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldTyp := field.Type
		origFieldTyp := fieldTyp

		if fieldTyp.Kind() == reflect.Ptr ||
			fieldTyp.Kind() == reflect.Slice ||
			fieldTyp.Kind() == reflect.Array ||
			fieldTyp.Kind() == reflect.Map {
			fieldTyp = fieldTyp.Elem()
		}
		if fieldTyp.Kind() != reflect.Struct {
			continue
		}

		var newPath string
		switch origFieldTyp.Kind() {
		case reflect.Struct, reflect.Ptr:
			newPath = fmt.Sprintf("%s.%s", path, field.Name)
		case reflect.Array, reflect.Slice:
			newPath = fmt.Sprintf("%s.%s[0]", path, field.Name)
		case reflect.Map:
			newPath = fmt.Sprintf("%s.%s[*]", path, field.Name)
		default:
			continue
		}

		tag := field.Tag.Get("json")
		name := parseTag(tag)
		if name == "" && !field.Anonymous {
			*errs = append(
				*errs,
				fmt.Errorf("field %s has no json tag and is not embedded, will cause unpredictable deserialization", newPath),
			)
		}
		checkType(fieldTyp, newPath, seenTypes, errs)
	}

	seenTypes[typ] = struct{}{}
}
