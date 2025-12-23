# dagger.io/dagger

```go
package dagger // import "dagger.io/dagger"

var EngineConn embed.FS
var GoDagGen []byte
var GoMod []byte
var GoSDK embed.FS
var GoSum []byte
var QueryBuilder embed.FS
var Telemetry embed.FS
func Connect(ctx context.Context, opts ...ClientOpt) (*Client, error)
func Tracer() trace.Tracer
type Address struct{ ... }
type AddressDirectoryOpts struct{ ... }
type AddressFileOpts struct{ ... }
type AddressID string
type Binding struct{ ... }
type BindingID string
type BuildArg struct{ ... }
type CacheSharingMode string
    const CacheSharingModeShared CacheSharingMode = "SHARED" ...
type CacheVolume struct{ ... }
type CacheVolumeID string
type Changeset struct{ ... }
type ChangesetID string
type Check struct{ ... }
type CheckGroup struct{ ... }
type CheckGroupID string
type CheckID string
type Client struct{ ... }
type ClientOpt interface{ ... }
    func WithConn(conn engineconn.EngineConn) ClientOpt
    func WithEnvironmentVariable(key, value string) ClientOpt
    func WithLogOutput(writer io.Writer) ClientOpt
    func WithRunnerHost(runnerHost string) ClientOpt
    func WithVerbosity(level int) ClientOpt
    func WithVersionOverride(version string) ClientOpt
    func WithWorkdir(path string) ClientOpt
type Cloud struct{ ... }
type CloudID string
type Container struct{ ... }
type ContainerAsServiceOpts struct{ ... }
type ContainerAsTarballOpts struct{ ... }
type ContainerDirectoryOpts struct{ ... }
type ContainerExistsOpts struct{ ... }
type ContainerExportImageOpts struct{ ... }
type ContainerExportOpts struct{ ... }
type ContainerFileOpts struct{ ... }
type ContainerID string
type ContainerImportOpts struct{ ... }
type ContainerOpts struct{ ... }
type ContainerPublishOpts struct{ ... }
type ContainerTerminalOpts struct{ ... }
type ContainerUpOpts struct{ ... }
type ContainerWithDefaultTerminalCmdOpts struct{ ... }
type ContainerWithDirectoryOpts struct{ ... }
type ContainerWithEntrypointOpts struct{ ... }
type ContainerWithEnvVariableOpts struct{ ... }
type ContainerWithExecOpts struct{ ... }
type ContainerWithExposedPortOpts struct{ ... }
type ContainerWithFileOpts struct{ ... }
type ContainerWithFilesOpts struct{ ... }
type ContainerWithMountedCacheOpts struct{ ... }
type ContainerWithMountedDirectoryOpts struct{ ... }
type ContainerWithMountedFileOpts struct{ ... }
type ContainerWithMountedSecretOpts struct{ ... }
type ContainerWithMountedTempOpts struct{ ... }
type ContainerWithNewFileOpts struct{ ... }
type ContainerWithSymlinkOpts struct{ ... }
type ContainerWithUnixSocketOpts struct{ ... }
type ContainerWithWorkdirOpts struct{ ... }
type ContainerWithoutDirectoryOpts struct{ ... }
type ContainerWithoutEntrypointOpts struct{ ... }
type ContainerWithoutExposedPortOpts struct{ ... }
type ContainerWithoutFileOpts struct{ ... }
type ContainerWithoutFilesOpts struct{ ... }
type ContainerWithoutMountOpts struct{ ... }
type ContainerWithoutUnixSocketOpts struct{ ... }
type CurrentModule struct{ ... }
type CurrentModuleID string
type CurrentModuleWorkdirOpts struct{ ... }
type DaggerObject = querybuilder.GraphQLMarshaller
type Directory struct{ ... }
type DirectoryAsModuleOpts struct{ ... }
type DirectoryAsModuleSourceOpts struct{ ... }
type DirectoryDockerBuildOpts struct{ ... }
type DirectoryEntriesOpts struct{ ... }
type DirectoryExistsOpts struct{ ... }
type DirectoryExportOpts struct{ ... }
type DirectoryFilterOpts struct{ ... }
type DirectoryID string
type DirectorySearchOpts struct{ ... }
type DirectoryTerminalOpts struct{ ... }
type DirectoryWithDirectoryOpts struct{ ... }
type DirectoryWithFileOpts struct{ ... }
type DirectoryWithFilesOpts struct{ ... }
type DirectoryWithNewDirectoryOpts struct{ ... }
type DirectoryWithNewFileOpts struct{ ... }
type Engine struct{ ... }
type EngineCache struct{ ... }
type EngineCacheEntry struct{ ... }
type EngineCacheEntryID string
type EngineCacheEntrySet struct{ ... }
type EngineCacheEntrySetID string
type EngineCacheEntrySetOpts struct{ ... }
type EngineCacheID string
type EngineCachePruneOpts struct{ ... }
type EngineID string
type EnumTypeDef struct{ ... }
type EnumTypeDefID string
type EnumValueTypeDef struct{ ... }
type EnumValueTypeDefID string
type Env struct{ ... }
type EnvFile struct{ ... }
type EnvFileGetOpts struct{ ... }
type EnvFileID string
type EnvFileOpts struct{ ... }
type EnvFileVariablesOpts struct{ ... }
type EnvID string
type EnvOpts struct{ ... }
type EnvVariable struct{ ... }
type EnvVariableID string
type Error struct{ ... }
type ErrorID string
type ErrorValue struct{ ... }
type ErrorValueID string
type ExecError struct{ ... }
type ExistsType string
    const ExistsTypeRegularType ExistsType = "REGULAR_TYPE" ...
type FieldTypeDef struct{ ... }
type FieldTypeDefID string
type File struct{ ... }
type FileAsEnvFileOpts struct{ ... }
type FileContentsOpts struct{ ... }
type FileDigestOpts struct{ ... }
type FileExportOpts struct{ ... }
type FileID string
type FileOpts struct{ ... }
type FileSearchOpts struct{ ... }
type FileWithReplacedOpts struct{ ... }
type Function struct{ ... }
type FunctionArg struct{ ... }
type FunctionArgID string
type FunctionCachePolicy string
    const FunctionCachePolicyDefault FunctionCachePolicy = "Default" ...
type FunctionCall struct{ ... }
type FunctionCallArgValue struct{ ... }
type FunctionCallArgValueID string
type FunctionCallID string
type FunctionID string
type FunctionWithArgOpts struct{ ... }
type FunctionWithCachePolicyOpts struct{ ... }
type FunctionWithDeprecatedOpts struct{ ... }
type GeneratedCode struct{ ... }
type GeneratedCodeID string
type GitOpts struct{ ... }
type GitRef struct{ ... }
type GitRefID string
type GitRefTreeOpts struct{ ... }
type GitRepository struct{ ... }
type GitRepositoryBranchesOpts struct{ ... }
type GitRepositoryID string
type GitRepositoryTagsOpts struct{ ... }
type HTTPOpts struct{ ... }
type Host struct{ ... }
type HostDirectoryOpts struct{ ... }
type HostFileOpts struct{ ... }
type HostFindUpOpts struct{ ... }
type HostID string
type HostServiceOpts struct{ ... }
type HostTunnelOpts struct{ ... }
type ImageLayerCompression string
    const ImageLayerCompressionGzip ImageLayerCompression = "Gzip" ...
type ImageMediaTypes string
    const ImageMediaTypesOcimediaTypes ImageMediaTypes = "OCIMediaTypes" ...
type InputTypeDef struct{ ... }
type InputTypeDefID string
type InterfaceTypeDef struct{ ... }
type InterfaceTypeDefID string
type JSON string
type JSONValue struct{ ... }
type JSONValueContentsOpts struct{ ... }
type JSONValueID string
type LLM struct{ ... }
type LLMID string
type LLMOpts struct{ ... }
type LLMTokenUsage struct{ ... }
type LLMTokenUsageID string
type Label struct{ ... }
type LabelID string
type ListTypeDef struct{ ... }
type ListTypeDefID string
type Module struct{ ... }
type ModuleChecksOpts struct{ ... }
type ModuleConfigClient struct{ ... }
type ModuleConfigClientID string
type ModuleID string
type ModuleServeOpts struct{ ... }
type ModuleSource struct{ ... }
type ModuleSourceExperimentalFeature string
    const ModuleSourceExperimentalFeatureSelfCalls ModuleSourceExperimentalFeature = "SELF_CALLS"
type ModuleSourceID string
type ModuleSourceKind string
    const ModuleSourceKindLocalSource ModuleSourceKind = "LOCAL_SOURCE" ...
type ModuleSourceOpts struct{ ... }
type NetworkProtocol string
    const NetworkProtocolTcp NetworkProtocol = "TCP" ...
type ObjectTypeDef struct{ ... }
type ObjectTypeDefID string
type PipelineLabel struct{ ... }
type Platform string
type Port struct{ ... }
type PortForward struct{ ... }
type PortID string
type Request struct{ ... }
type Response struct{ ... }
type ReturnType string
    const ReturnTypeSuccess ReturnType = "SUCCESS" ...
type SDKConfig struct{ ... }
type SDKConfigID string
type ScalarTypeDef struct{ ... }
type ScalarTypeDefID string
type SearchResult struct{ ... }
type SearchResultID string
type SearchSubmatch struct{ ... }
type SearchSubmatchID string
type Secret struct{ ... }
type SecretID string
type SecretOpts struct{ ... }
type Service struct{ ... }
type ServiceEndpointOpts struct{ ... }
type ServiceID string
type ServiceStopOpts struct{ ... }
type ServiceTerminalOpts struct{ ... }
type ServiceUpOpts struct{ ... }
type Socket struct{ ... }
type SocketID string
type SourceMap struct{ ... }
type SourceMapID string
type Terminal struct{ ... }
type TerminalID string
type TypeDef struct{ ... }
type TypeDefID string
type TypeDefKind string
    const TypeDefKindStringKind TypeDefKind = "STRING_KIND" ...
type TypeDefWithEnumMemberOpts struct{ ... }
type TypeDefWithEnumOpts struct{ ... }
type TypeDefWithEnumValueOpts struct{ ... }
type TypeDefWithFieldOpts struct{ ... }
type TypeDefWithInterfaceOpts struct{ ... }
type TypeDefWithObjectOpts struct{ ... }
type TypeDefWithScalarOpts struct{ ... }
type Void string
type WithCheckFunc func(r *Check) *Check
type WithCheckGroupFunc func(r *CheckGroup) *CheckGroup
type WithContainerFunc func(r *Container) *Container
type WithDirectoryFunc func(r *Directory) *Directory
type WithEnvFileFunc func(r *EnvFile) *EnvFile
type WithEnvFunc func(r *Env) *Env
type WithErrorFunc func(r *Error) *Error
type WithFileFunc func(r *File) *File
type WithFunctionFunc func(r *Function) *Function
type WithGeneratedCodeFunc func(r *GeneratedCode) *GeneratedCode
type WithGitRefFunc func(r *GitRef) *GitRef
type WithJSONValueFunc func(r *JSONValue) *JSONValue
type WithLLMFunc func(r *LLM) *LLM
type WithModuleFunc func(r *Module) *Module
type WithModuleSourceFunc func(r *ModuleSource) *ModuleSource
type WithServiceFunc func(r *Service) *Service
type WithTypeDefFunc func(r *TypeDef) *TypeDef
```

