// Code generated by pluginator on PrefixSuffixTransformer; DO NOT EDIT.
// pluginator {Version:unknown GitCommit:$Format:%H$ BuildDate:1970-01-01T00:00:00Z GoOs:linux GoArch:amd64}

package builtins

import (
	"errors"
	"fmt"

	"sigs.k8s.io/kustomize/api/transform"
	"sigs.k8s.io/kustomize/api/types"

	"sigs.k8s.io/kustomize/api/resid"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/yaml"
)

// Add the given prefix and suffix to the field.
type PrefixSuffixTransformerPlugin struct {
	Prefix     string            `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Suffix     string            `json:"suffix,omitempty" yaml:"suffix,omitempty"`
	FieldSpecs []types.FieldSpec `json:"fieldSpecs,omitempty" yaml:"fieldSpecs,omitempty"`
}

// Not placed in a file yet due to lack of demand.
var prefixSuffixFieldSpecsToSkip = []types.FieldSpec{
	{
		Gvk: resid.Gvk{Kind: "CustomResourceDefinition"},
	},
	{
		Gvk: resid.Gvk{Group: "apiregistration.k8s.io", Kind: "APIService"},
	},
}

func (p *PrefixSuffixTransformerPlugin) Config(
	h *resmap.PluginHelpers, c []byte) (err error) {
	p.Prefix = ""
	p.Suffix = ""
	p.FieldSpecs = nil
	err = yaml.Unmarshal(c, p)
	if err != nil {
		return
	}
	if p.FieldSpecs == nil {
		return errors.New("fieldSpecs is not expected to be nil")
	}
	return
}

func (p *PrefixSuffixTransformerPlugin) Transform(m resmap.ResMap) error {

	// Even if both the Prefix and Suffix are empty we want
	// to proceed with the transformation. This allows to add contextual
	// information to the resources (AddNamePrefix and AddNameSuffix).

	for _, r := range m.Resources() {
		if p.shouldSkip(r.OrgId()) {
			// Don't change the actual definition
			// of a CRD.
			continue
		}
		id := r.OrgId()
		// current default configuration contains
		// only one entry: "metadata/name" with no GVK
		for _, path := range p.FieldSpecs {
			if !id.IsSelected(&path.Gvk) {
				// With the currrent default configuration,
				// because no Gvk is specified, so a wild
				// card
				continue
			}

			if smellsLikeANameChange(&path) {
				// "metadata/name" is the only field.
				// this will add a prefix and a suffix
				// to the resource even if those are
				// empty
				r.AddNamePrefix(p.Prefix)
				r.AddNameSuffix(p.Suffix)
			}

			// the addPrefixSuffix method will not
			// change the name if both the prefix and suffix
			// are empty.
			err := transform.MutateField(
				r.Map(),
				path.PathSlice(),
				path.CreateIfNotPresent,
				p.addPrefixSuffix)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func smellsLikeANameChange(fs *types.FieldSpec) bool {
	return fs.Path == "metadata/name"
}

func (p *PrefixSuffixTransformerPlugin) shouldSkip(
	id resid.ResId) bool {
	for _, path := range prefixSuffixFieldSpecsToSkip {
		if id.IsSelected(&path.Gvk) {
			return true
		}
	}
	return false
}

func (p *PrefixSuffixTransformerPlugin) addPrefixSuffix(
	in interface{}) (interface{}, error) {
	s, ok := in.(string)
	if !ok {
		return nil, fmt.Errorf("%#v is expected to be %T", in, s)
	}
	return fmt.Sprintf("%s%s%s", p.Prefix, s, p.Suffix), nil
}

func NewPrefixSuffixTransformerPlugin() resmap.TransformerPlugin {
	return &PrefixSuffixTransformerPlugin{}
}
