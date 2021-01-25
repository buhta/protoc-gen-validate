package ocaml

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/envoyproxy/protoc-gen-validate/templates/shared"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/iancoleman/strcase"
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func RegisterIndex(tpl *template.Template, params pgs.Parameters) {
	fns := ocamlFuncs{pgsgo.InitContext(params)}

	tpl.Funcs(map[string]interface{}{
		"classNameFile": classNameFile,
		"importsPvg":    importsPvg,
		"ocamlPackage":  ocamlPackage,
		"simpleName":    fns.Name,
		"qualifiedName": fns.qualifiedName,
	})
}

func Register(tpl *template.Template, params pgs.Parameters) {
	fns := ocamlFuncs{pgsgo.InitContext(params)}

	tpl.Funcs(map[string]interface{}{
		"accessor":                  fns.accessor,
		"byteArrayLit":              fns.byteArrayLit,
		"camelCase":                 fns.camelCase,
		"classNameFile":             classNameFile,
		"classNameMessage":          classNameMessage,
		"durLit":                    fns.durLit,
		"fieldName":                 fns.fieldName,
		"ocamlPackage":              ocamlPackage,
		"ocamlStringEscape":         fns.ocamlStringEscape,
		"ocamlTypeFor":              fns.ocamlTypeFor,
		"ocamlTypeLiteralSuffixFor": fns.ocamlTypeLiteralSuffixFor,
		"hasAccessor":               fns.hasAccessor,
		"oneof":                     fns.oneofTypeName,
		"sprintf":                   fmt.Sprintf,
		"simpleName":                fns.Name,
		"tsLit":                     fns.tsLit,
		"qualifiedName":             fns.qualifiedName,
		"isOfFileType":              fns.isOfFileType,
		"isOfMessageType":           fns.isOfMessageType,
		"isOfStringType":            fns.isOfStringType,
		"unwrap":                    fns.unwrap,
		"renderConstants":           fns.renderConstants(tpl),
		"constantName":              fns.constantName,
	})

	template.Must(tpl.Parse(fileTpl))
	template.Must(tpl.New("msg").Parse(msgTpl))
	template.Must(tpl.New("msgInner").Parse(msgInnerTpl))

	template.Must(tpl.New("wrapper").Parse(wrapperTpl))
	template.Must(tpl.New("wrapperConst").Parse(wrapperConstTpl))
}

type ocamlFuncs struct{ pgsgo.Context }

func PhpFilePath(f pgs.File, ctx pgsgo.Context, tpl *template.Template) *pgs.FilePath {
	// Don't generate validators for files that don't import PGV
	if !importsPvg(f) {
		return nil
	}

	fullPath := strings.Replace(ocamlPackage(f), ".", string(os.PathSeparator), -1)
	fileName := classNameFile(f) + "Validator.ocaml"
	filePath := pgs.JoinPaths(fullPath, fileName)
	return &filePath
}

func JavaMultiFilePath(f pgs.File, m pgs.Message) pgs.FilePath {
	fullPath := strings.Replace(ocamlPackage(f), ".", string(os.PathSeparator), -1)
	fileName := classNameMessage(m) + "Validator.ocaml"
	filePath := pgs.JoinPaths(fullPath, fileName)
	return filePath
}

func importsPvg(f pgs.File) bool {
	for _, dep := range f.Descriptor().Dependency {
		if strings.HasSuffix(dep, "validate.proto") {
			return true
		}
	}
	return false
}

func classNameFile(f pgs.File) string {
	// Explicit outer class name overrides implicit name
	options := f.Descriptor().GetOptions()
	if options != nil && !options.GetJavaMultipleFiles() && options.JavaOuterClassname != nil {
		return options.GetJavaOuterClassname()
	}

	protoName := pgs.FilePath(f.Name().String()).BaseName()

	className := sanitizeClassName(protoName)
	className = appendOuterClassName(className, f)

	return className
}

func classNameMessage(m pgs.Message) string {
	className := m.Name().String()
	// This is really silly, but when the multiple files option is true, protoc puts underscores in file names.
	// When multiple files is false, underscores are stripped. Short of rewriting all the name sanitization
	// logic for ocaml, using "UnderscoreUnderscoreUnderscore" is an escape sequence seems to work with an extremely
	// small likelihood of name conflict.
	className = strings.Replace(className, "_", "UnderscoreUnderscoreUnderscore", -1)
	className = sanitizeClassName(className)
	className = strings.Replace(className, "UnderscoreUnderscoreUnderscore", "_", -1)
	return className
}

func sanitizeClassName(className string) string {
	className = makeInvalidClassnameCharactersUnderscores(className)
	className = underscoreBetweenConsecutiveUppercase(className)
	className = strcase.ToCamel(strcase.ToSnake(className))
	className = upperCaseAfterNumber(className)
	return className
}

func ocamlPackage(file pgs.File) string {
	// Explicit ocaml package overrides implicit package
	options := file.Descriptor().GetOptions()
	if options != nil && options.JavaPackage != nil {
		return options.GetJavaPackage()
	}
	return strcase.ToCamel(file.Package().ProtoName().String())
}

func (fns ocamlFuncs) qualifiedName(entity pgs.Entity) string {
	file, isFile := entity.(pgs.File)
	if isFile {
		name := ocamlPackage(file)
		if file.Descriptor().GetOptions() != nil {
			if !file.Descriptor().GetOptions().GetJavaMultipleFiles() {
				name += ("." + classNameFile(file))
			}
		} else {
			name += ("." + classNameFile(file))
		}
		return name
	}

	message, isMessage := entity.(pgs.Message)
	if isMessage && message.Parent() != nil {
		// recurse
		return fns.qualifiedName(message.Parent()) + "." + entity.Name().String()
	}

	enum, isEnum := entity.(pgs.Enum)
	if isEnum && enum.Parent() != nil {
		// recurse
		return fns.qualifiedName(enum.Parent()) + "." + entity.Name().String()
	}

	return entity.Name().String()
}

// Replace invalid identifier characters with an underscore
func makeInvalidClassnameCharactersUnderscores(name string) string {
	var sb string
	for _, c := range name {
		switch {
		case c >= '0' && c <= '9':
			sb += string(c)
		case c >= 'a' && c <= 'z':
			sb += string(c)
		case c >= 'A' && c <= 'Z':
			sb += string(c)
		default:
			sb += "_"
		}
	}
	return sb
}

func upperCaseAfterNumber(name string) string {
	var sb string
	var p rune

	for _, c := range name {
		if unicode.IsDigit(p) {
			sb += string(unicode.ToUpper(c))
		} else {
			sb += string(c)
		}
		p = c
	}
	return sb
}

func underscoreBetweenConsecutiveUppercase(name string) string {
	var sb string
	var p rune

	for _, c := range name {
		if unicode.IsUpper(p) && unicode.IsUpper(c) {
			sb += "_" + string(c)
		} else {
			sb += string(c)
		}
		p = c
	}
	return sb
}

func appendOuterClassName(outerClassName string, file pgs.File) string {
	conflict := false

	for _, enum := range file.Enums() {
		if enum.Name().String() == outerClassName {
			conflict = true
		}
	}

	for _, message := range file.Messages() {
		if message.Name().String() == outerClassName {
			conflict = true
		}
	}

	for _, service := range file.Services() {
		if service.Name().String() == outerClassName {
			conflict = true
		}
	}

	if conflict {
		return outerClassName + "OuterClass"
	} else {
		return outerClassName
	}
}

func (fns ocamlFuncs) accessor(ctx shared.RuleContext) string {
	if ctx.AccessorOverride != "" {
		return ctx.AccessorOverride
	}
	return fns.fieldAccessor(ctx.Field)
}

func (fns ocamlFuncs) fieldAccessor(f pgs.Field) string {
	fieldName := strcase.ToCamel(f.Name().String())
	if f.Type().IsMap() {
		fieldName += "Map"
	}
	if f.Type().IsRepeated() {
		fieldName += "List"
	}

	fieldName = upperCaseAfterNumber(fieldName)
	return fmt.Sprintf("proto.get%s()", fieldName)
}

func (fns ocamlFuncs) hasAccessor(ctx shared.RuleContext) string {
	if ctx.AccessorOverride != "" {
		return "true"
	}
	fiedlName := strcase.ToCamel(ctx.Field.Name().String())
	fiedlName = upperCaseAfterNumber(fiedlName)
	return "proto.has" + fiedlName + "()"
}

func (fns ocamlFuncs) fieldName(ctx shared.RuleContext) string {
	return ctx.Field.Name().String()
}

func (fns ocamlFuncs) ocamlTypeFor(ctx shared.RuleContext) string {
	t := ctx.Field.Type()

	// Map key and value types
	if t.IsMap() {
		switch ctx.AccessorOverride {
		case "key":
			return fns.ocamlTypeForProtoType(t.Key().ProtoType())
		case "value":
			return fns.ocamlTypeForProtoType(t.Element().ProtoType())
		}
	}

	if t.IsEmbed() {
		if embed := t.Embed(); embed.IsWellKnown() {
			switch embed.WellKnownType() {
			case pgs.AnyWKT:
				return "String"
			case pgs.DurationWKT:
				return "com.google.protobuf.Duration"
			case pgs.TimestampWKT:
				return "com.google.protobuf.Timestamp"
			case pgs.Int32ValueWKT, pgs.UInt32ValueWKT:
				return "Integer"
			case pgs.Int64ValueWKT, pgs.UInt64ValueWKT:
				return "Long"
			case pgs.DoubleValueWKT:
				return "Double"
			case pgs.FloatValueWKT:
				return "Float"
			}
		}
	}

	if t.IsRepeated() {
		if t.ProtoType() == pgs.MessageT {
			return fns.qualifiedName(t.Element().Embed())
		} else if t.ProtoType() == pgs.EnumT {
			return fns.qualifiedName(t.Element().Enum())
		}
	}

	if t.IsEnum() {
		return fns.qualifiedName(t.Enum())
	}

	return fns.ocamlTypeForProtoType(t.ProtoType())
}

func (fns ocamlFuncs) ocamlTypeForProtoType(t pgs.ProtoType) string {

	switch t {
	case pgs.Int32T, pgs.UInt32T, pgs.SInt32, pgs.Fixed32T, pgs.SFixed32:
		return "Integer"
	case pgs.Int64T, pgs.UInt64T, pgs.SInt64, pgs.Fixed64T, pgs.SFixed64:
		return "Long"
	case pgs.DoubleT:
		return "Double"
	case pgs.FloatT:
		return "Float"
	case pgs.BoolT:
		return "Boolean"
	case pgs.StringT:
		return "String"
	case pgs.BytesT:
		return "com.google.protobuf.ByteString"
	default:
		return "Object"
	}
}

func (fns ocamlFuncs) ocamlTypeLiteralSuffixFor(ctx shared.RuleContext) string {
	t := ctx.Field.Type()

	if t.IsMap() {
		switch ctx.AccessorOverride {
		case "key":
			return fns.ocamlTypeLiteralSuffixForPrototype(t.Key().ProtoType())
		case "value":
			return fns.ocamlTypeLiteralSuffixForPrototype(t.Element().ProtoType())
		}
	}

	if t.IsEmbed() {
		if embed := t.Embed(); embed.IsWellKnown() {
			switch embed.WellKnownType() {
			case pgs.Int64ValueWKT, pgs.UInt64ValueWKT:
				return "L"
			case pgs.FloatValueWKT:
				return "F"
			case pgs.DoubleValueWKT:
				return "D"
			}
		}
	}

	return fns.ocamlTypeLiteralSuffixForPrototype(t.ProtoType())
}

func (fns ocamlFuncs) ocamlTypeLiteralSuffixForPrototype(t pgs.ProtoType) string {
	switch t {
	case pgs.Int64T, pgs.UInt64T, pgs.SInt64, pgs.Fixed64T, pgs.SFixed64:
		return "L"
	case pgs.FloatT:
		return "F"
	case pgs.DoubleT:
		return "D"
	default:
		return ""
	}
}

func (fns ocamlFuncs) ocamlStringEscape(s string) string {
	s = fmt.Sprintf("%q", s)
	s = s[1 : len(s)-1]
	s = strings.Replace(s, `\u00`, `\x`, -1)
	s = strings.Replace(s, `\x`, `\\x`, -1)
	// s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, `"`, `\"`, -1)
	return `"` + s + `"`
}

func (fns ocamlFuncs) camelCase(name pgs.Name) string {
	return strcase.ToCamel(name.String())
}

func (fns ocamlFuncs) byteArrayLit(bytes []uint8) string {
	var sb string
	sb += "new byte[]{"
	for _, b := range bytes {
		sb += fmt.Sprintf("(byte)%#x,", b)
	}
	sb += "}"

	return sb
}

func (fns ocamlFuncs) durLit(dur *duration.Duration) string {
	return fmt.Sprintf(
		"io.envoyproxy.pgv.TimestampValidation.toDuration(%d,%d)",
		dur.GetSeconds(), dur.GetNanos())
}

func (fns ocamlFuncs) tsLit(ts *timestamp.Timestamp) string {
	return fmt.Sprintf(
		"io.envoyproxy.pgv.TimestampValidation.toTimestamp(%d,%d)",
		ts.GetSeconds(), ts.GetNanos())
}

func (fns ocamlFuncs) oneofTypeName(f pgs.Field) pgsgo.TypeName {
	return pgsgo.TypeName(fmt.Sprintf("%s", strings.ToUpper(f.Name().String())))
}

func (fns ocamlFuncs) isOfFileType(o interface{}) bool {
	switch o.(type) {
	case pgs.File:
		return true
	default:
		return false
	}
}

func (fns ocamlFuncs) isOfMessageType(f pgs.Field) bool {
	return f.Type().ProtoType() == pgs.MessageT
}

func (fns ocamlFuncs) isOfStringType(f pgs.Field) bool {
	return f.Type().ProtoType() == pgs.StringT
}

func (fns ocamlFuncs) unwrap(ctx shared.RuleContext) (shared.RuleContext, error) {
	ctx, err := ctx.Unwrap("wrapped")
	if err != nil {
		return ctx, err
	}
	ctx.AccessorOverride = fmt.Sprintf("%s.get%s()", fns.fieldAccessor(ctx.Field),
		fns.camelCase(ctx.Field.Type().Embed().Fields()[0].Name()))
	return ctx, nil
}

func (fns ocamlFuncs) renderConstants(tpl *template.Template) func(ctx shared.RuleContext) (string, error) {
	return func(ctx shared.RuleContext) (string, error) {
		var b bytes.Buffer
		var err error

		hasConstTemplate := false
		for _, t := range tpl.Templates() {
			if t.Name() == ctx.Typ+"Const" {
				hasConstTemplate = true
			}
		}

		if hasConstTemplate {
			err = tpl.ExecuteTemplate(&b, ctx.Typ+"Const", ctx)
		}

		return b.String(), err
	}
}

func (fns ocamlFuncs) constantName(ctx shared.RuleContext, rule string) string {
	return strcase.ToScreamingSnake(ctx.Field.Name().String() + "_" + ctx.Index + "_" + rule)
}
