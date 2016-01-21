package dartlang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"gopkg.in/yaml.v2"

	"github.com/Workiva/frugal/compiler/generator"
	"github.com/Workiva/frugal/compiler/globals"
	"github.com/Workiva/frugal/compiler/parser"
)

const (
	lang               = "dart"
	defaultOutputDir   = "gen-dart"
	serviceSuffix      = "_service"
	scopeSuffix        = "_scope"
	minimumDartVersion = "1.12.0"
	tab                = "  "
	tabtab             = tab + tab
	tabtabtab          = tab + tab + tab
	tabtabtabtab       = tab + tab + tab + tab
	tabtabtabtabtab    = tab + tab + tab + tab + tab
)

type Generator struct {
	*generator.BaseGenerator
}

func NewGenerator(options map[string]string) generator.LanguageGenerator {
	return &Generator{&generator.BaseGenerator{Options: options}}
}

func (g *Generator) GetOutputDir(dir string) string {
	if pkg, ok := g.Frugal.Thrift.Namespace(lang); ok {
		dir = filepath.Join(dir, toLibraryName(pkg))
	} else {
		dir = filepath.Join(dir, g.Frugal.Name)
	}
	return dir
}

func (g *Generator) DefaultOutputDir() string {
	return defaultOutputDir
}

func (g *Generator) GenerateDependencies(dir string) error {
	if err := g.addToPubspec(dir); err != nil {
		return err
	}
	if err := g.exportClasses(dir); err != nil {
		return err
	}
	return nil
}

type pubspec struct {
	Name         string                      `yaml:"name"`
	Version      string                      `yaml:"version"`
	Description  string                      `yaml:"description"`
	Environment  env                         `yaml:"environment"`
	Dependencies map[interface{}]interface{} `yaml:"dependencies"`
}

type env struct {
	SDK string `yaml:"sdk"`
}

type dep struct {
	Hosted  hostedDep `yaml:"hosted,omitempty"`
	Path    string    `yaml:"path,omitempty"`
	Version string    `yaml:"version,omitempty"`
}

type hostedDep struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func (g *Generator) addToPubspec(dir string) error {
	pubFilePath := filepath.Join(dir, "pubspec.yaml")

	deps := map[interface{}]interface{}{
		"thrift": dep{Hosted: hostedDep{Name: "thrift", URL: "https://pub.workiva.org"}, Version: "^0.0.1"},
	}

	if g.Frugal.ContainsFrugalDefinitions() {
		deps["frugal"] = dep{Hosted: hostedDep{Name: "frugal", URL: "https://pub.workiva.org"}, Version: "^0.0.1"}
	}

	for _, include := range g.Frugal.ReferencedIncludes() {
		namespace, ok := g.Frugal.NamespaceForInclude(include, lang)
		if !ok {
			namespace = include
		}
		deps[toLibraryName(namespace)] = dep{Path: "../" + toLibraryName(namespace)}
	}

	namespace, ok := g.Frugal.Thrift.Namespace(lang)
	if !ok {
		namespace = g.Frugal.Name
	}

	ps := &pubspec{
		Name:        strings.ToLower(toLibraryName(namespace)),
		Version:     globals.Version,
		Description: "Autogenerated by the frugal compiler",
		Environment: env{
			SDK: "^" + minimumDartVersion,
		},
		Dependencies: deps,
	}

	d, err := yaml.Marshal(&ps)
	if err != nil {
		return err
	}
	// create and write to new file
	newPubFile, err := os.Create(pubFilePath)
	defer newPubFile.Close()
	if err != nil {
		return err
	}
	if _, err := newPubFile.Write(d); err != nil {
		return err
	}
	return nil
}

func (g *Generator) exportClasses(dir string) error {
	filename := strings.ToLower(g.Frugal.Name)
	if ns, ok := g.Frugal.Thrift.Namespace(lang); ok {
		filename = strings.ToLower(toLibraryName(ns))
	}
	dartFile := fmt.Sprintf("%s.%s", filename, lang)
	mainFilePath := filepath.Join(dir, "lib", dartFile)
	mainFile, err := os.OpenFile(mainFilePath, syscall.O_RDWR, 0777)
	defer mainFile.Close()
	if err != nil {
		return err
	}

	exports := "\n"
	for _, service := range g.Frugal.Thrift.Services {
		servTitle := strings.Title(service.Name)
		exports += fmt.Sprintf("export 'src/%s%s%s.%s' show F%sClient;\n",
			generator.FilePrefix, strings.ToLower(service.Name), serviceSuffix, lang, servTitle)
	}
	for _, scope := range g.Frugal.Scopes {
		scopeTitle := strings.Title(scope.Name)
		exports += fmt.Sprintf("export 'src/%s%s%s.%s' show %sPublisher, %sSubscriber;\n",
			generator.FilePrefix, strings.ToLower(scope.Name), scopeSuffix, lang, scopeTitle, scopeTitle)
	}
	stat, err := mainFile.Stat()
	if err != nil {
		return err
	}
	_, err = mainFile.WriteAt([]byte(exports), stat.Size())
	return err
}

func (g *Generator) GenerateFile(name, outputDir string, fileType generator.FileType) (*os.File, error) {
	outputDir = filepath.Join(outputDir, "lib")
	outputDir = filepath.Join(outputDir, "src")
	switch fileType {
	case generator.CombinedServiceFile:
		return g.CreateFile(strings.ToLower(name)+serviceSuffix, outputDir, lang, true)
	case generator.CombinedScopeFile:
		return g.CreateFile(strings.ToLower(name)+scopeSuffix, outputDir, lang, true)
	default:
		return nil, fmt.Errorf("frugal: Bad file type for dartlang generator: %s", fileType)
	}
}

func (g *Generator) GenerateDocStringComment(file *os.File) error {
	comment := fmt.Sprintf(
		"// Autogenerated by Frugal Compiler (%s)\n"+
			"// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING",
		globals.Version)

	_, err := file.WriteString(comment)
	return err
}

func (g *Generator) GenerateServicePackage(file *os.File, s *parser.Service) error {
	return g.generatePackage(file, s.Name, serviceSuffix)
}

func (g *Generator) GenerateScopePackage(file *os.File, s *parser.Scope) error {
	return g.generatePackage(file, s.Name, scopeSuffix)
}

func (g *Generator) generatePackage(file *os.File, name, suffix string) error {
	pkg, ok := g.Frugal.Thrift.Namespace(lang)
	if ok {
		components := generator.GetPackageComponents(pkg)
		pkg = components[len(components)-1]
	} else {
		pkg = g.Frugal.Name
	}
	_, err := file.WriteString(fmt.Sprintf("library %s.src.%s%s%s;", pkg,
		generator.FilePrefix, strings.ToLower(name), scopeSuffix))
	return err
}

func (g *Generator) GenerateServiceImports(file *os.File, s *parser.Service) error {
	imports := "import 'dart:async';\n\n"
	imports += "import 'dart:typed_data' show Uint8List;\n"
	imports += "import 'package:thrift/thrift.dart' as thrift;\n"
	imports += "import 'package:frugal/frugal.dart' as frugal;\n\n"
	// import included packages
	for _, include := range s.ReferencedIncludes() {
		namespace, ok := g.Frugal.Thrift.NamespaceForInclude(include, lang)
		if !ok {
			namespace = include
		}
		namespace = strings.ToLower(toLibraryName(namespace))
		imports += fmt.Sprintf("import 'package:%s/%s.dart' as t_%s;\n", namespace, namespace, namespace)
	}

	// Import same package.
	pkgLower := strings.ToLower(g.Frugal.Name)
	imports += fmt.Sprintf("import 'package:%s/%s.dart' as t_%s;\n", pkgLower, pkgLower, pkgLower)

	// Import thrift package for method args
	servLower := strings.ToLower(s.Name)
	imports += fmt.Sprintf("import '%s.dart' as t_%s;\n", servLower, servLower)

	_, err := file.WriteString(imports)
	return err
}

func (g *Generator) GenerateScopeImports(file *os.File, s *parser.Scope) error {
	imports := "import 'dart:async';\n\n"
	imports += "import 'package:thrift/thrift.dart' as thrift;\n"
	imports += "import 'package:frugal/frugal.dart' as frugal;\n\n"
	// import included packages
	for _, include := range s.ReferencedIncludes() {
		namespace, ok := s.Frugal.NamespaceForInclude(include, lang)
		if !ok {
			namespace = include
		}
		namespace = strings.ToLower(toLibraryName(namespace))
		imports += fmt.Sprintf("import 'package:%s/%s.dart' as t_%s;\n", namespace, namespace, namespace)
	}

	// Import same package.
	pkgLower := strings.ToLower(g.Frugal.Name)
	imports += fmt.Sprintf("import 'package:%s/%s.dart' as t_%s;\n", pkgLower, pkgLower, pkgLower)

	_, err := file.WriteString(imports)
	return err
}

func (g *Generator) GenerateConstants(file *os.File, name string) error {
	constants := fmt.Sprintf("const String delimiter = '%s';", globals.TopicDelimiter)
	_, err := file.WriteString(constants)
	return err
}

func (g *Generator) GeneratePublisher(file *os.File, scope *parser.Scope) error {
	publishers := ""
	if scope.Comment != nil {
		publishers += g.GenerateInlineComment(scope.Comment, "/")
	}
	publishers += fmt.Sprintf("class %sPublisher {\n", strings.Title(scope.Name))
	publishers += tab + "frugal.FScopeTransport fTransport;\n"
	publishers += tab + "frugal.FProtocol fProtocol;\n"

	publishers += fmt.Sprintf(tab+"%sPublisher(frugal.FScopeProvider provider) {\n", strings.Title(scope.Name))
	publishers += tabtab + "fTransport = provider.fTransportFactory.getTransport();\n"
	publishers += tabtab + "fProtocol = provider.fProtocolFactory.getProtocol(fTransport);\n"
	publishers += tab + "}\n\n"

	publishers += tab + "Future open() {\n"
	publishers += tabtab + "return fTransport.open();\n"
	publishers += tab + "}\n\n"

	publishers += tab + "Future close() {\n"
	publishers += tabtab + "return fTransport.close();\n"
	publishers += tab + "}\n\n"

	args := ""
	if len(scope.Prefix.Variables) > 0 {
		for _, variable := range scope.Prefix.Variables {
			args = fmt.Sprintf("%sString %s, ", args, variable)
		}
	}
	prefix := ""
	for _, op := range scope.Operations {
		publishers += prefix
		prefix = "\n\n"
		if op.Comment != nil {
			publishers += g.GenerateInlineComment(op.Comment, tab+"/")
		}
		publishers += fmt.Sprintf(tab+"Future publish%s(frugal.FContext ctx, %s%s req) async {\n", op.Name, args, g.qualifiedParamName(op))
		publishers += fmt.Sprintf(tabtab+"var op = \"%s\";\n", op.Name)
		publishers += fmt.Sprintf(tabtab+"var prefix = \"%s\";\n", generatePrefixStringTemplate(scope))
		publishers += tabtab + "var topic = \"${prefix}" + strings.Title(scope.Name) + "${delimiter}${op}\";\n"
		publishers += tabtab + "fTransport.setTopic(topic);\n"
		publishers += tabtab + "var oprot = fProtocol;\n"
		publishers += tabtab + "var msg = new thrift.TMessage(op, thrift.TMessageType.CALL, 0);\n"
		publishers += tabtab + "oprot.writeRequestHeader(ctx);\n"
		publishers += tabtab + "oprot.writeMessageBegin(msg);\n"
		publishers += tabtab + "req.write(oprot);\n"
		publishers += tabtab + "oprot.writeMessageEnd();\n"
		publishers += tabtab + "await oprot.transport.flush();\n"
		publishers += tab + "}\n"
	}

	publishers += "}\n"

	_, err := file.WriteString(publishers)
	return err
}

func generatePrefixStringTemplate(scope *parser.Scope) string {
	if scope.Prefix.String == "" {
		return ""
	}
	template := ""
	template += scope.Prefix.Template()
	template += globals.TopicDelimiter
	if len(scope.Prefix.Variables) == 0 {
		return template
	}
	vars := make([]interface{}, len(scope.Prefix.Variables))
	for i, variable := range scope.Prefix.Variables {
		vars[i] = fmt.Sprintf("${%s}", variable)
	}
	template = fmt.Sprintf(template, vars...)
	return template
}

func (g *Generator) GenerateSubscriber(file *os.File, scope *parser.Scope) error {
	subscribers := ""
	if scope.Comment != nil {
		subscribers += g.GenerateInlineComment(scope.Comment, "/")
	}
	subscribers += fmt.Sprintf("class %sSubscriber {\n", strings.Title(scope.Name))
	subscribers += tab + "final frugal.FScopeProvider provider;\n\n"

	subscribers += fmt.Sprintf(tab+"%sSubscriber(this.provider) {}\n\n", strings.Title(scope.Name))

	args := ""
	if len(scope.Prefix.Variables) > 0 {
		for _, variable := range scope.Prefix.Variables {
			args = fmt.Sprintf("%sString %s, ", args, variable)
		}
	}
	prefix := ""
	for _, op := range scope.Operations {
		subscribers += prefix
		prefix = "\n\n"
		if op.Comment != nil {
			subscribers += g.GenerateInlineComment(op.Comment, tab+"/")
		}
		subscribers += fmt.Sprintf(tab+"Future<frugal.FSubscription> subscribe%s(%sdynamic on%s(frugal.FContext ctx, %s req)) async {\n",
			op.Name, args, op.ParamName(), g.qualifiedParamName(op))
		subscribers += fmt.Sprintf(tabtab+"var op = \"%s\";\n", op.Name)
		subscribers += fmt.Sprintf(tabtab+"var prefix = \"%s\";\n", generatePrefixStringTemplate(scope))
		subscribers += tabtab + "var topic = \"${prefix}" + strings.Title(scope.Name) + "${delimiter}${op}\";\n"
		subscribers += tabtab + "var transport = provider.fTransportFactory.getTransport();\n"
		subscribers += fmt.Sprintf(tabtab+"await transport.subscribe(topic, _recv%s(op, provider.fProtocolFactory, on%s));\n",
			op.Name, op.ParamName())
		subscribers += tabtab + "return new frugal.FSubscription(topic, transport);\n"
		subscribers += tab + "}\n\n"

		subscribers += fmt.Sprintf(tab+"_recv%s(String op, frugal.FProtocolFactory protocolFactory, dynamic on%s(frugal.FContext ctx, %s req)) {\n",
			op.Name, op.ParamName(), g.qualifiedParamName(op))
		subscribers += fmt.Sprintf(tabtab+"callback%s(thrift.TTransport transport) {\n", op.Name)
		subscribers += tabtabtab + "var iprot = protocolFactory.getProtocol(transport);\n"
		subscribers += tabtabtab + "var ctx = iprot.readRequestHeader();\n"
		subscribers += tabtabtab + "var tMsg = iprot.readMessageBegin();\n"
		subscribers += tabtabtab + "if (tMsg.name != op) {\n"
		subscribers += tabtabtabtab + "thrift.TProtocolUtil.skip(iprot, thrift.TType.STRUCT);\n"
		subscribers += tabtabtabtab + "iprot.readMessageEnd();\n"
		subscribers += tabtabtabtab + "throw new thrift.TApplicationError(\n"
		subscribers += tabtabtabtab + "thrift.TApplicationErrorType.UNKNOWN_METHOD, tMsg.name);\n"
		subscribers += tabtabtab + "}\n"
		subscribers += fmt.Sprintf(tabtabtab+"var req = new %s();\n", g.qualifiedParamName(op))
		subscribers += tabtabtab + "req.read(iprot);\n"
		subscribers += tabtabtab + "iprot.readMessageEnd();\n"
		subscribers += fmt.Sprintf(tabtabtab+"on%s(ctx, req);\n", op.ParamName())
		subscribers += tabtab + "}\n"
		subscribers += fmt.Sprintf(tabtab+"return callback%s;\n", op.Name)
		subscribers += tab + "}\n"
	}

	subscribers += "}\n"

	_, err := file.WriteString(subscribers)
	return err
}

func (g *Generator) GenerateService(file *os.File, s *parser.Service) error {
	contents := ""
	contents += g.generateInterface(s)
	contents += g.generateClient(s)

	_, err := file.WriteString(contents)
	return err
}

func (g *Generator) generateInterface(service *parser.Service) string {
	contents := ""
	if service.Comment != nil {
		contents += g.GenerateInlineComment(service.Comment, "/")
	}
	contents += fmt.Sprintf("abstract class F%s {\n", strings.Title(service.Name))
	for _, method := range service.Methods {
		contents += "\n"
		if method.Comment != nil {
			contents += g.GenerateInlineComment(method.Comment, tab+"/")
		}
		contents += fmt.Sprintf(tab+"Future%s %s(frugal.FContext ctx%s);\n",
			g.generateReturnArg(method), strings.ToLower(method.Name), g.generateInputArgs(method.Arguments))
	}
	contents += "}\n\n"
	return contents
}

func (g *Generator) generateClient(service *parser.Service) string {
	servTitle := strings.Title(service.Name)
	contents := ""
	if service.Comment != nil {
		contents += g.GenerateInlineComment(service.Comment, "/")
	}
	contents += fmt.Sprintf("class F%sClient implements F%s {\n", servTitle, servTitle)
	contents += "\n"
	contents += fmt.Sprintf(tab+"F%sClient(frugal.FServiceProvider provider) {\n", servTitle)
	contents += tabtab + "_transport = provider.fTransport;\n"
	contents += tabtab + "_transport.setRegistry(new frugal.FClientRegistry());\n"
	contents += tabtab + "_protocolFactory = provider.fProtocolFactory;\n"
	contents += tabtab + "_oprot = _protocolFactory.getProtocol(_transport);\n"
	contents += tab + "}\n\n"
	contents += tab + "frugal.FTransport _transport;\n"
	contents += tab + "frugal.FProtocolFactory _protocolFactory;\n"
	contents += tab + "frugal.FProtocol _oprot;\n"
	contents += tab + "frugal.FProtocol get oprot => _oprot;\n\n"

	for _, method := range service.Methods {
		contents += g.generateClientMethod(service, method)
	}
	contents += "}\n"
	return contents
}

func (g *Generator) generateClientMethod(service *parser.Service, method *parser.Method) string {
	servLower := strings.ToLower(service.Name)
	nameTitle := strings.Title(method.Name)
	nameLower := strings.ToLower(method.Name)

	contents := ""
	if method.Comment != nil {
		contents += g.GenerateInlineComment(method.Comment, tab+"/")
	}
	// Generate the calling method
	contents += fmt.Sprintf(tab+"Future%s %s(frugal.FContext ctx%s) async {\n",
		g.generateReturnArg(method), nameLower, g.generateInputArgs(method.Arguments))
	contents += tabtab + "var controller = new StreamController();\n"
	contents += fmt.Sprintf(tabtab+"_transport.register(ctx, _recv%sHandler(ctx, controller));\n", nameTitle)
	contents += tabtab + "try {\n"
	contents += tabtabtab + "oprot.writeRequestHeader(ctx);\n"
	contents += fmt.Sprintf(tabtabtab+"oprot.writeMessageBegin(new thrift.TMessage(\"%s\", thrift.TMessageType.CALL, 0));\n",
		nameLower)
	contents += fmt.Sprintf(tabtabtab+"t_%s.%s_args args = new t_%s.%s_args();\n",
		servLower, nameLower, servLower, nameLower)
	for _, arg := range method.Arguments {
		argLower := strings.ToLower(arg.Name)
		contents += fmt.Sprintf(tabtabtab+"args.%s = %s;\n", argLower, argLower)
	}
	contents += tabtabtab + "args.write(oprot);\n"
	contents += tabtabtab + "oprot.writeMessageEnd();\n"
	contents += tabtabtab + "await oprot.transport.flush();\n"
	contents += tabtabtab + "return await controller.stream.first.timeout(ctx.timeout);\n"
	contents += tabtab + "} finally {\n"
	contents += tabtabtab + "_transport.unregister(ctx);\n"
	contents += tabtab + "}\n"
	contents += tab + "}\n\n"

	// Generate the callback
	contents += fmt.Sprintf(tab+"_recv%sHandler(frugal.FContext ctx, StreamController controller) {\n", nameTitle)
	contents += fmt.Sprintf(tabtab+"%sCallback(thrift.TTransport transport) {\n", nameLower)
	contents += tabtabtab + "try {\n"
	contents += tabtabtabtab + "var iprot = _protocolFactory.getProtocol(transport);\n"
	contents += tabtabtabtab + "iprot.readResponseHeader(ctx);\n"
	contents += tabtabtabtab + "thrift.TMessage msg = iprot.readMessageBegin();\n"
	contents += tabtabtabtab + "if (msg.type == thrift.TMessageType.EXCEPTION) {\n"
	contents += tabtabtabtabtab + "thrift.TApplicationError error = thrift.TApplicationError.read(iprot);\n"
	contents += tabtabtabtabtab + "iprot.readMessageEnd();\n"
	contents += tabtabtabtabtab + "throw error;\n"
	contents += tabtabtabtab + "}\n\n"

	contents += fmt.Sprintf(tabtabtabtab+"t_%s.%s_result result = new t_%s.%s_result();\n",
		servLower, nameLower, servLower, nameLower)
	contents += tabtabtabtab + "result.read(iprot);\n"
	contents += tabtabtabtab + "iprot.readMessageEnd();\n"
	if method.ReturnType == nil {
		contents += g.generateErrors(method)
		contents += tabtabtabtab + "controller.add(null);\n"
	} else {
		contents += tabtabtabtab + "if (result.isSetSuccess()) {\n"
		contents += tabtabtabtabtab + "controller.add(result.success);\n"
		contents += tabtabtabtabtab + "return;\n"
		contents += tabtabtabtab + "}\n\n"
		contents += g.generateErrors(method)
		contents += tabtabtabtab + "throw new thrift.TApplicationError(\n"
		contents += fmt.Sprintf(tabtabtabtabtab+"thrift.TApplicationErrorType.MISSING_RESULT, "+
			"\"%s failed: unknown result\"\n",
			nameLower)
		contents += tabtabtabtab + ");\n"
	}
	contents += tabtabtab + "} catch(e) {\n"
	contents += tabtabtabtab + "controller.addError(e);\n"
	contents += tabtabtabtab + "rethrow;\n"
	contents += tabtabtab + "}\n"
	contents += tabtab + "}\n"
	contents += fmt.Sprintf(tabtab+"return %sCallback;\n", nameLower)
	contents += tab + "}\n\n"

	return contents
}

func (g *Generator) generateReturnArg(method *parser.Method) string {
	if method.ReturnType == nil {
		return ""
	}
	return fmt.Sprintf("<%s>", g.getDartTypeFromThriftType(method.ReturnType))
}

func (g *Generator) generateInputArgs(args []*parser.Field) string {
	argStr := ""
	for _, arg := range args {
		argStr += ", " + g.getDartTypeFromThriftType(arg.Type) + " " + strings.ToLower(arg.Name)
	}
	return argStr
}

func (g *Generator) generateErrors(method *parser.Method) string {
	contents := ""
	for _, exp := range method.Exceptions {
		contents += fmt.Sprintf(tabtabtabtab+"if (result.%s != null) {\n", strings.ToLower(exp.Name))
		contents += fmt.Sprintf(tabtabtabtabtab+"controller.addError(result.%s);\n", strings.ToLower(exp.Name))
		contents += tabtabtabtabtab + "return;\n"
		contents += tabtabtabtab + "}\n"
	}
	return contents
}

func (g *Generator) getDartTypeFromThriftType(t *parser.Type) string {
	if t == nil {
		return "void"
	}
	underlyingType := g.Frugal.UnderlyingType(t)
	switch underlyingType.Name {
	case "bool":
		return "bool"
	case "byte":
		return "int"
	case "i16":
		return "int"
	case "i32":
		return "int"
	case "i64":
		return "int"
	case "double":
		return "double"
	case "string":
		return "String"
	case "binary":
		return "Uint8List"
	case "list":
		return fmt.Sprintf("List<%s>",
			g.getDartTypeFromThriftType(underlyingType.ValueType))
	case "set":
		return fmt.Sprintf("Set<%s>",
			g.getDartTypeFromThriftType(underlyingType.ValueType))
	case "map":
		return fmt.Sprintf("Map<%s,%s>",
			g.getDartTypeFromThriftType(underlyingType.KeyType),
			g.getDartTypeFromThriftType(underlyingType.ValueType))
	default:
		// This is a custom type
		return g.qualifiedTypeName(t)
	}
}

// get qualafied type names for custom types
func (g *Generator) qualifiedTypeName(t *parser.Type) string {
	param := strings.Title(t.ParamName())
	include := t.IncludeName()
	if include != "" {
		namespace, ok := g.Frugal.NamespaceForInclude(include, lang)
		if !ok {
			namespace = include
		}
		namespace = toLibraryName(namespace)
		param = fmt.Sprintf("t_%s.%s", strings.ToLower(namespace), param)
	} else {
		param = fmt.Sprintf("t_%s.%s", strings.ToLower(g.Frugal.Name), param)
	}
	return param
}

func (g *Generator) qualifiedParamName(op *parser.Operation) string {
	param := op.ParamName()
	include := op.IncludeName()
	if include != "" {
		namespace, ok := g.Frugal.NamespaceForInclude(include, lang)
		if !ok {
			namespace = include
		}
		namespace = toLibraryName(namespace)
		param = fmt.Sprintf("t_%s.%s", strings.ToLower(namespace), param)
	} else {
		param = fmt.Sprintf("t_%s.%s", strings.ToLower(param), param)
	}
	return param
}

func toLibraryName(name string) string {
	return strings.Replace(name, ".", "_", -1)
}
