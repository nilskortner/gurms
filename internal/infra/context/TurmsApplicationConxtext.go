package context

const BUILD_INFO_PROPS_PATH = "git.properties"
const PROPERTY_BUILD_VERSION = "git.build.version"
const PROPERTY_BUILD_TIME = "git.build.time"
const PROPERTY_COMMIT_ID = "git.commit.id.full"
const DEFAULT_VERSION = "0.0.0"

var isClosing bool

var home string
var configDir string
var tempDir string

var isProduction bool
var idDevorLocaltest bool
var activeEnvProfile string
var buildProperties BuildProperties

var shutdownJobTimeoutMillis int64
var shutdownHooks

func TurmsApplicationContext() {

}