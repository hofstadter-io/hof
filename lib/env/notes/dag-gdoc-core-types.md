# dagger.Client

```go
package dagger // import "dagger.io/dagger"

type Client struct {
	// Has unexported fields.
}
    Client is the Dagger Engine Client

func (r *Client) Address(value string) *Address
func (r *Client) CacheVolume(key string) *CacheVolume
func (c *Client) Close() error
func (r *Client) Cloud() *Cloud
func (r *Client) Container(opts ...ContainerOpts) *Container
func (r *Client) CurrentEnv() *Env
func (r *Client) CurrentFunctionCall() *FunctionCall
func (r *Client) CurrentModule() *CurrentModule
func (r *Client) CurrentTypeDefs(ctx context.Context) ([]TypeDef, error)
func (r *Client) DefaultPlatform(ctx context.Context) (Platform, error)
func (r *Client) Directory() *Directory
func (c *Client) Do(ctx context.Context, req *Request, resp *Response) error
func (r *Client) Engine() *Engine
func (r *Client) Env(opts ...EnvOpts) *Env
func (r *Client) EnvFile(opts ...EnvFileOpts) *EnvFile
func (r *Client) Error(message string) *Error
func (r *Client) File(name string, contents string, opts ...FileOpts) *File
func (r *Client) Function(name string, returnType *TypeDef) *Function
func (r *Client) GeneratedCode(code *Directory) *GeneratedCode
func (r *Client) Git(url string, opts ...GitOpts) *GitRepository
func (c *Client) GraphQLClient() graphql.Client
func (r *Client) HTTP(url string, opts ...HTTPOpts) *File
func (r *Client) Host() *Host
func (r *Client) JSON() *JSONValue
func (r *Client) LLM(opts ...LLMOpts) *LLM
func (r *Client) LoadAddressFromID(id AddressID) *Address
func (r *Client) LoadBindingFromID(id BindingID) *Binding
func (r *Client) LoadCacheVolumeFromID(id CacheVolumeID) *CacheVolume
func (r *Client) LoadChangesetFromID(id ChangesetID) *Changeset
func (r *Client) LoadCheckFromID(id CheckID) *Check
func (r *Client) LoadCheckGroupFromID(id CheckGroupID) *CheckGroup
func (r *Client) LoadCloudFromID(id CloudID) *Cloud
func (r *Client) LoadContainerFromID(id ContainerID) *Container
func (r *Client) LoadCurrentModuleFromID(id CurrentModuleID) *CurrentModule
func (r *Client) LoadDirectoryFromID(id DirectoryID) *Directory
func (r *Client) LoadEngineCacheEntryFromID(id EngineCacheEntryID) *EngineCacheEntry
func (r *Client) LoadEngineCacheEntrySetFromID(id EngineCacheEntrySetID) *EngineCacheEntrySet
func (r *Client) LoadEngineCacheFromID(id EngineCacheID) *EngineCache
func (r *Client) LoadEngineFromID(id EngineID) *Engine
func (r *Client) LoadEnumTypeDefFromID(id EnumTypeDefID) *EnumTypeDef
func (r *Client) LoadEnumValueTypeDefFromID(id EnumValueTypeDefID) *EnumValueTypeDef
func (r *Client) LoadEnvFileFromID(id EnvFileID) *EnvFile
func (r *Client) LoadEnvFromID(id EnvID) *Env
func (r *Client) LoadEnvVariableFromID(id EnvVariableID) *EnvVariable
func (r *Client) LoadErrorFromID(id ErrorID) *Error
func (r *Client) LoadErrorValueFromID(id ErrorValueID) *ErrorValue
func (r *Client) LoadFieldTypeDefFromID(id FieldTypeDefID) *FieldTypeDef
func (r *Client) LoadFileFromID(id FileID) *File
func (r *Client) LoadFunctionArgFromID(id FunctionArgID) *FunctionArg
func (r *Client) LoadFunctionCallArgValueFromID(id FunctionCallArgValueID) *FunctionCallArgValue
func (r *Client) LoadFunctionCallFromID(id FunctionCallID) *FunctionCall
func (r *Client) LoadFunctionFromID(id FunctionID) *Function
func (r *Client) LoadGeneratedCodeFromID(id GeneratedCodeID) *GeneratedCode
func (r *Client) LoadGitRefFromID(id GitRefID) *GitRef
func (r *Client) LoadGitRepositoryFromID(id GitRepositoryID) *GitRepository
func (r *Client) LoadHostFromID(id HostID) *Host
func (r *Client) LoadInputTypeDefFromID(id InputTypeDefID) *InputTypeDef
func (r *Client) LoadInterfaceTypeDefFromID(id InterfaceTypeDefID) *InterfaceTypeDef
func (r *Client) LoadJSONValueFromID(id JSONValueID) *JSONValue
func (r *Client) LoadLLMFromID(id LLMID) *LLM
func (r *Client) LoadLLMTokenUsageFromID(id LLMTokenUsageID) *LLMTokenUsage
func (r *Client) LoadLabelFromID(id LabelID) *Label
func (r *Client) LoadListTypeDefFromID(id ListTypeDefID) *ListTypeDef
func (r *Client) LoadModuleConfigClientFromID(id ModuleConfigClientID) *ModuleConfigClient
func (r *Client) LoadModuleFromID(id ModuleID) *Module
func (r *Client) LoadModuleSourceFromID(id ModuleSourceID) *ModuleSource
func (r *Client) LoadObjectTypeDefFromID(id ObjectTypeDefID) *ObjectTypeDef
func (r *Client) LoadPortFromID(id PortID) *Port
func (r *Client) LoadSDKConfigFromID(id SDKConfigID) *SDKConfig
func (r *Client) LoadScalarTypeDefFromID(id ScalarTypeDefID) *ScalarTypeDef
func (r *Client) LoadSearchResultFromID(id SearchResultID) *SearchResult
func (r *Client) LoadSearchSubmatchFromID(id SearchSubmatchID) *SearchSubmatch
func (r *Client) LoadSecretFromID(id SecretID) *Secret
func (r *Client) LoadServiceFromID(id ServiceID) *Service
func (r *Client) LoadSocketFromID(id SocketID) *Socket
func (r *Client) LoadSourceMapFromID(id SourceMapID) *SourceMap
func (r *Client) LoadTerminalFromID(id TerminalID) *Terminal
func (r *Client) LoadTypeDefFromID(id TypeDefID) *TypeDef
func (r *Client) Module() *Module
func (r *Client) ModuleSource(refString string, opts ...ModuleSourceOpts) *ModuleSource
func (c *Client) QueryBuilder() *querybuilder.Selection
func (r *Client) Secret(uri string, opts ...SecretOpts) *Secret
func (r *Client) SetSecret(name string, plaintext string) *Secret
func (r *Client) SourceMap(filename string, line int, column int) *SourceMap
func (r *Client) TypeDef() *TypeDef
func (r *Client) Version(ctx context.Context) (string, error)
func (r *Client) WithGraphQLQuery(q *querybuilder.Selection) *Client
```


# dagger.Container

```go
package dagger // import "dagger.io/dagger"

type Container struct {
	// Has unexported fields.
}
    An OCI-compatible container, also known as a Docker container.

func (r *Container) AsService(opts ...ContainerAsServiceOpts) *Service
func (r *Container) AsTarball(opts ...ContainerAsTarballOpts) *File
func (r *Container) CombinedOutput(ctx context.Context) (string, error)
func (r *Container) DefaultArgs(ctx context.Context) ([]string, error)
func (r *Container) Directory(path string, opts ...ContainerDirectoryOpts) *Directory
func (r *Container) Entrypoint(ctx context.Context) ([]string, error)
func (r *Container) EnvVariable(ctx context.Context, name string) (string, error)
func (r *Container) EnvVariables(ctx context.Context) ([]EnvVariable, error)
func (r *Container) Exists(ctx context.Context, path string, opts ...ContainerExistsOpts) (bool, error)
func (r *Container) ExitCode(ctx context.Context) (int, error)
func (r *Container) ExperimentalWithAllGPUs() *Container
func (r *Container) ExperimentalWithGPU(devices []string) *Container
func (r *Container) Export(ctx context.Context, path string, opts ...ContainerExportOpts) (string, error)
func (r *Container) ExportImage(ctx context.Context, name string, opts ...ContainerExportImageOpts) error
func (r *Container) ExposedPorts(ctx context.Context) ([]Port, error)
func (r *Container) File(path string, opts ...ContainerFileOpts) *File
func (r *Container) From(address string) *Container
func (r *Container) ID(ctx context.Context) (ContainerID, error)
func (r *Container) ImageRef(ctx context.Context) (string, error)
func (r *Container) Import(source *File, opts ...ContainerImportOpts) *Container
func (r *Container) Label(ctx context.Context, name string) (string, error)
func (r *Container) Labels(ctx context.Context) ([]Label, error)
func (r *Container) MarshalJSON() ([]byte, error)
func (r *Container) Mounts(ctx context.Context) ([]string, error)
func (r *Container) Platform(ctx context.Context) (Platform, error)
func (r *Container) Publish(ctx context.Context, address string, opts ...ContainerPublishOpts) (string, error)
func (r *Container) Rootfs() *Directory
func (r *Container) Stderr(ctx context.Context) (string, error)
func (r *Container) Stdout(ctx context.Context) (string, error)
func (r *Container) Sync(ctx context.Context) (*Container, error)
func (r *Container) Terminal(opts ...ContainerTerminalOpts) *Container
func (r *Container) Up(ctx context.Context, opts ...ContainerUpOpts) error
func (r *Container) User(ctx context.Context) (string, error)
func (r *Container) With(f WithContainerFunc) *Container
func (r *Container) WithAnnotation(name string, value string) *Container
func (r *Container) WithDefaultArgs(args []string) *Container
func (r *Container) WithDefaultTerminalCmd(args []string, opts ...ContainerWithDefaultTerminalCmdOpts) *Container
func (r *Container) WithDirectory(path string, source *Directory, opts ...ContainerWithDirectoryOpts) *Container
func (r *Container) WithEntrypoint(args []string, opts ...ContainerWithEntrypointOpts) *Container
func (r *Container) WithEnvVariable(name string, value string, opts ...ContainerWithEnvVariableOpts) *Container
func (r *Container) WithError(err string) *Container
func (r *Container) WithExec(args []string, opts ...ContainerWithExecOpts) *Container
func (r *Container) WithExposedPort(port int, opts ...ContainerWithExposedPortOpts) *Container
func (r *Container) WithFile(path string, source *File, opts ...ContainerWithFileOpts) *Container
func (r *Container) WithFiles(path string, sources []*File, opts ...ContainerWithFilesOpts) *Container
func (r *Container) WithGraphQLQuery(q *querybuilder.Selection) *Container
func (r *Container) WithLabel(name string, value string) *Container
func (r *Container) WithMountedCache(path string, cache *CacheVolume, opts ...ContainerWithMountedCacheOpts) *Container
func (r *Container) WithMountedDirectory(path string, source *Directory, opts ...ContainerWithMountedDirectoryOpts) *Container
func (r *Container) WithMountedFile(path string, source *File, opts ...ContainerWithMountedFileOpts) *Container
func (r *Container) WithMountedSecret(path string, source *Secret, opts ...ContainerWithMountedSecretOpts) *Container
func (r *Container) WithMountedTemp(path string, opts ...ContainerWithMountedTempOpts) *Container
func (r *Container) WithNewFile(path string, contents string, opts ...ContainerWithNewFileOpts) *Container
func (r *Container) WithRegistryAuth(address string, username string, secret *Secret) *Container
func (r *Container) WithRootfs(directory *Directory) *Container
func (r *Container) WithSecretVariable(name string, secret *Secret) *Container
func (r *Container) WithServiceBinding(alias string, service *Service) *Container
func (r *Container) WithSymlink(target string, linkName string, opts ...ContainerWithSymlinkOpts) *Container
func (r *Container) WithUnixSocket(path string, source *Socket, opts ...ContainerWithUnixSocketOpts) *Container
func (r *Container) WithUser(name string) *Container
func (r *Container) WithWorkdir(path string, opts ...ContainerWithWorkdirOpts) *Container
func (r *Container) WithoutAnnotation(name string) *Container
func (r *Container) WithoutDefaultArgs() *Container
func (r *Container) WithoutDirectory(path string, opts ...ContainerWithoutDirectoryOpts) *Container
func (r *Container) WithoutEntrypoint(opts ...ContainerWithoutEntrypointOpts) *Container
func (r *Container) WithoutEnvVariable(name string) *Container
func (r *Container) WithoutExposedPort(port int, opts ...ContainerWithoutExposedPortOpts) *Container
func (r *Container) WithoutFile(path string, opts ...ContainerWithoutFileOpts) *Container
func (r *Container) WithoutFiles(paths []string, opts ...ContainerWithoutFilesOpts) *Container
func (r *Container) WithoutLabel(name string) *Container
func (r *Container) WithoutMount(path string, opts ...ContainerWithoutMountOpts) *Container
func (r *Container) WithoutRegistryAuth(address string) *Container
func (r *Container) WithoutSecretVariable(name string) *Container
func (r *Container) WithoutUnixSocket(path string, opts ...ContainerWithoutUnixSocketOpts) *Container
func (r *Container) WithoutUser() *Container
func (r *Container) WithoutWorkdir() *Container
func (r *Container) Workdir(ctx context.Context) (string, error)
func (r *Container) XXX_GraphQLID(ctx context.Context) (string, error)
func (r *Container) XXX_GraphQLIDType() string
func (r *Container) XXX_GraphQLType() string
```


# dagger.Directory

```go
package dagger // import "dagger.io/dagger"

type Directory struct {
	// Has unexported fields.
}
    A directory.

func (r *Directory) AsGit() *GitRepository
func (r *Directory) AsModule(opts ...DirectoryAsModuleOpts) *Module
func (r *Directory) AsModuleSource(opts ...DirectoryAsModuleSourceOpts) *ModuleSource
func (r *Directory) Changes(from *Directory) *Changeset
func (r *Directory) Chown(path string, owner string) *Directory
func (r *Directory) Diff(other *Directory) *Directory
func (r *Directory) Digest(ctx context.Context) (string, error)
func (r *Directory) Directory(path string) *Directory
func (r *Directory) DockerBuild(opts ...DirectoryDockerBuildOpts) *Container
func (r *Directory) Entries(ctx context.Context, opts ...DirectoryEntriesOpts) ([]string, error)
func (r *Directory) Exists(ctx context.Context, path string, opts ...DirectoryExistsOpts) (bool, error)
func (r *Directory) Export(ctx context.Context, path string, opts ...DirectoryExportOpts) (string, error)
func (r *Directory) File(path string) *File
func (r *Directory) Filter(opts ...DirectoryFilterOpts) *Directory
func (r *Directory) FindUp(ctx context.Context, name string, start string) (string, error)
func (r *Directory) Glob(ctx context.Context, pattern string) ([]string, error)
func (r *Directory) ID(ctx context.Context) (DirectoryID, error)
func (r *Directory) MarshalJSON() ([]byte, error)
func (r *Directory) Name(ctx context.Context) (string, error)
func (r *Directory) Search(ctx context.Context, pattern string, opts ...DirectorySearchOpts) ([]SearchResult, error)
func (r *Directory) Sync(ctx context.Context) (*Directory, error)
func (r *Directory) Terminal(opts ...DirectoryTerminalOpts) *Directory
func (r *Directory) With(f WithDirectoryFunc) *Directory
func (r *Directory) WithChanges(changes *Changeset) *Directory
func (r *Directory) WithDirectory(path string, source *Directory, opts ...DirectoryWithDirectoryOpts) *Directory
func (r *Directory) WithError(err string) *Directory
func (r *Directory) WithFile(path string, source *File, opts ...DirectoryWithFileOpts) *Directory
func (r *Directory) WithFiles(path string, sources []*File, opts ...DirectoryWithFilesOpts) *Directory
func (r *Directory) WithGraphQLQuery(q *querybuilder.Selection) *Directory
func (r *Directory) WithNewDirectory(path string, opts ...DirectoryWithNewDirectoryOpts) *Directory
func (r *Directory) WithNewFile(path string, contents string, opts ...DirectoryWithNewFileOpts) *Directory
func (r *Directory) WithPatch(patch string) *Directory
func (r *Directory) WithPatchFile(patch *File) *Directory
func (r *Directory) WithSymlink(target string, linkName string) *Directory
func (r *Directory) WithTimestamps(timestamp int) *Directory
func (r *Directory) WithoutDirectory(path string) *Directory
func (r *Directory) WithoutFile(path string) *Directory
func (r *Directory) WithoutFiles(paths []string) *Directory
func (r *Directory) XXX_GraphQLID(ctx context.Context) (string, error)
func (r *Directory) XXX_GraphQLIDType() string
func (r *Directory) XXX_GraphQLType() string
```


# dagger.File

```go
package dagger // import "dagger.io/dagger"

type File struct {
	// Has unexported fields.
}
    A file.

func (r *File) AsEnvFile(opts ...FileAsEnvFileOpts) *EnvFile
func (r *File) Chown(owner string) *File
func (r *File) Contents(ctx context.Context, opts ...FileContentsOpts) (string, error)
func (r *File) Digest(ctx context.Context, opts ...FileDigestOpts) (string, error)
func (r *File) Export(ctx context.Context, path string, opts ...FileExportOpts) (string, error)
func (r *File) ID(ctx context.Context) (FileID, error)
func (r *File) MarshalJSON() ([]byte, error)
func (r *File) Name(ctx context.Context) (string, error)
func (r *File) Search(ctx context.Context, pattern string, opts ...FileSearchOpts) ([]SearchResult, error)
func (r *File) Size(ctx context.Context) (int, error)
func (r *File) Sync(ctx context.Context) (*File, error)
func (r *File) With(f WithFileFunc) *File
func (r *File) WithGraphQLQuery(q *querybuilder.Selection) *File
func (r *File) WithName(name string) *File
func (r *File) WithReplaced(search string, replacement string, opts ...FileWithReplacedOpts) *File
func (r *File) WithTimestamps(timestamp int) *File
func (r *File) XXX_GraphQLID(ctx context.Context) (string, error)
func (r *File) XXX_GraphQLIDType() string
func (r *File) XXX_GraphQLType() string
```


# dagger.Service

```go
package dagger // import "dagger.io/dagger"

type Service struct {
	// Has unexported fields.
}
    A content-addressed service providing TCP connectivity.

func (r *Service) Endpoint(ctx context.Context, opts ...ServiceEndpointOpts) (string, error)
func (r *Service) Hostname(ctx context.Context) (string, error)
func (r *Service) ID(ctx context.Context) (ServiceID, error)
func (r *Service) MarshalJSON() ([]byte, error)
func (r *Service) Ports(ctx context.Context) ([]Port, error)
func (r *Service) Start(ctx context.Context) (*Service, error)
func (r *Service) Stop(ctx context.Context, opts ...ServiceStopOpts) (*Service, error)
func (r *Service) Sync(ctx context.Context) (*Service, error)
func (r *Service) Terminal(opts ...ServiceTerminalOpts) *Service
func (r *Service) Up(ctx context.Context, opts ...ServiceUpOpts) error
func (r *Service) With(f WithServiceFunc) *Service
func (r *Service) WithGraphQLQuery(q *querybuilder.Selection) *Service
func (r *Service) WithHostname(hostname string) *Service
func (r *Service) XXX_GraphQLID(ctx context.Context) (string, error)
func (r *Service) XXX_GraphQLIDType() string
func (r *Service) XXX_GraphQLType() string
```


