export namespace browser {
	
	export class Browser {
	    name: string;
	    id: string;
	    path: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Browser(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.id = source["id"];
	        this.path = source["path"];
	        this.enabled = source["enabled"];
	    }
	}
	export class CleanItem {
	    browserId: string;
	    browserName: string;
	    category: string;
	    label: string;
	    path: string;
	    fileCount: number;
	    totalSize: number;
	    checked: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CleanItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.browserId = source["browserId"];
	        this.browserName = source["browserName"];
	        this.category = source["category"];
	        this.label = source["label"];
	        this.path = source["path"];
	        this.fileCount = source["fileCount"];
	        this.totalSize = source["totalSize"];
	        this.checked = source["checked"];
	        this.error = source["error"];
	    }
	}
	export class CleanResult {
	    browserId: string;
	    category: string;
	    label: string;
	    success: boolean;
	    fileCount: number;
	    freedBytes: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CleanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.browserId = source["browserId"];
	        this.category = source["category"];
	        this.label = source["label"];
	        this.success = source["success"];
	        this.fileCount = source["fileCount"];
	        this.freedBytes = source["freedBytes"];
	        this.error = source["error"];
	    }
	}

}

export namespace clipboard {
	
	export class ClipItem {
	    id: number;
	    content: string;
	    type: string;
	    time: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new ClipItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.content = source["content"];
	        this.type = source["type"];
	        this.time = source["time"];
	        this.size = source["size"];
	    }
	}

}

export namespace common {
	
	export class AppConfig {
	    theme: string;
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.language = source["language"];
	    }
	}

}

export namespace devtools {
	
	export class CodeBeautifyResult {
	    input: string;
	    output: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CodeBeautifyResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input = source["input"];
	        this.output = source["output"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class CodecResult {
	    input: string;
	    output: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CodecResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input = source["input"];
	        this.output = source["output"];
	        this.error = source["error"];
	    }
	}
	export class ColorResult {
	    hex: string;
	    rgb: string;
	    hsl: string;
	    hsv: string;
	    name?: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ColorResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hex = source["hex"];
	        this.rgb = source["rgb"];
	        this.hsl = source["hsl"];
	        this.hsv = source["hsv"];
	        this.name = source["name"];
	        this.error = source["error"];
	    }
	}
	export class DiffResult {
	    type: string;
	    oldLine: string;
	    newLine: string;
	    oldNum: number;
	    newNum: number;
	
	    static createFrom(source: any = {}) {
	        return new DiffResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.oldLine = source["oldLine"];
	        this.newLine = source["newLine"];
	        this.oldNum = source["oldNum"];
	        this.newNum = source["newNum"];
	    }
	}
	export class FormatResult {
	    input: string;
	    output: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new FormatResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input = source["input"];
	        this.output = source["output"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class HTTPRequest {
	    url: string;
	    method: string;
	    headers: string;
	    body: string;
	
	    static createFrom(source: any = {}) {
	        return new HTTPRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.method = source["method"];
	        this.headers = source["headers"];
	        this.body = source["body"];
	    }
	}
	export class HTTPResponse {
	    statusCode: number;
	    statusText: string;
	    headers: string;
	    body: string;
	    duration: string;
	    size: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new HTTPResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.statusCode = source["statusCode"];
	        this.statusText = source["statusText"];
	        this.headers = source["headers"];
	        this.body = source["body"];
	        this.duration = source["duration"];
	        this.size = source["size"];
	        this.error = source["error"];
	    }
	}
	export class JSONResult {
	    formatted: string;
	    valid: boolean;
	    error?: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new JSONResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.formatted = source["formatted"];
	        this.valid = source["valid"];
	        this.error = source["error"];
	        this.size = source["size"];
	    }
	}
	export class JWTResult {
	    raw: string;
	    header: string;
	    payload: string;
	    signature: string;
	    valid: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new JWTResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.raw = source["raw"];
	        this.header = source["header"];
	        this.payload = source["payload"];
	        this.signature = source["signature"];
	        this.valid = source["valid"];
	        this.error = source["error"];
	    }
	}
	export class QRResult {
	    content: string;
	    dataUri: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new QRResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.dataUri = source["dataUri"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class RegexTestResult {
	    pattern: string;
	    text: string;
	    matches: string[];
	    count: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new RegexTestResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pattern = source["pattern"];
	        this.text = source["text"];
	        this.matches = source["matches"];
	        this.count = source["count"];
	        this.error = source["error"];
	    }
	}
	export class ReleaseInfo {
	    tag_name: string;
	    name: string;
	    body: string;
	    published_at: string;
	    html_url: string;
	    prerelease: boolean;
	    download_url?: string;
	
	    static createFrom(source: any = {}) {
	        return new ReleaseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag_name = source["tag_name"];
	        this.name = source["name"];
	        this.body = source["body"];
	        this.published_at = source["published_at"];
	        this.html_url = source["html_url"];
	        this.prerelease = source["prerelease"];
	        this.download_url = source["download_url"];
	    }
	}
	export class TimestampResult {
	    unixTimestamp: number;
	    dateTime: string;
	    iso8601: string;
	
	    static createFrom(source: any = {}) {
	        return new TimestampResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.unixTimestamp = source["unixTimestamp"];
	        this.dateTime = source["dateTime"];
	        this.iso8601 = source["iso8601"];
	    }
	}
	export class UUIDGenResult {
	    uuids: string[];
	    count: number;
	    version: number;
	
	    static createFrom(source: any = {}) {
	        return new UUIDGenResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuids = source["uuids"];
	        this.count = source["count"];
	        this.version = source["version"];
	    }
	}

}

export namespace filetools {
	
	export class BatchCompressResult {
	    inputPath: string;
	    outputPath: string;
	    originalSize: number;
	    newSize: number;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchCompressResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputPath = source["inputPath"];
	        this.outputPath = source["outputPath"];
	        this.originalSize = source["originalSize"];
	        this.newSize = source["newSize"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class ConvertResult {
	    inputPath: string;
	    outputPath: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ConvertResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputPath = source["inputPath"];
	        this.outputPath = source["outputPath"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class DiffLine {
	    type: string;
	    oldLine?: number;
	    newLine?: number;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new DiffLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.oldLine = source["oldLine"];
	        this.newLine = source["newLine"];
	        this.content = source["content"];
	    }
	}
	export class DocConvertResult {
	    inputPath: string;
	    outputPath: string;
	    targetType: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new DocConvertResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputPath = source["inputPath"];
	        this.outputPath = source["outputPath"];
	        this.targetType = source["targetType"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class FileMatch {
	    size: number;
	    modTime: string;
	    path: string;
	    hash: string;
	    matchType: string;
	
	    static createFrom(source: any = {}) {
	        return new FileMatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.size = source["size"];
	        this.modTime = source["modTime"];
	        this.path = source["path"];
	        this.hash = source["hash"];
	        this.matchType = source["matchType"];
	    }
	}
	export class DuplicateGroup {
	    hash: string;
	    fileCount: number;
	    totalSize: number;
	    files: FileMatch[];
	
	    static createFrom(source: any = {}) {
	        return new DuplicateGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	        this.fileCount = source["fileCount"];
	        this.totalSize = source["totalSize"];
	        this.files = this.convertValues(source["files"], FileMatch);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EncodingResult {
	    input: string;
	    output: string;
	    from: string;
	    to: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new EncodingResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.input = source["input"];
	        this.output = source["output"];
	        this.from = source["from"];
	        this.to = source["to"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	
	export class FolderSize {
	    path: string;
	    name: string;
	    size: number;
	    files: number;
	    pct: number;
	
	    static createFrom(source: any = {}) {
	        return new FolderSize(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.files = source["files"];
	        this.pct = source["pct"];
	    }
	}
	export class LargeFile {
	    path: string;
	    name: string;
	    size: number;
	    modified: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new LargeFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.modified = source["modified"];
	        this.type = source["type"];
	    }
	}
	export class OrganizePreview {
	    sourcePath: string;
	    destPath: string;
	    sourceName: string;
	    folderName: string;
	    fileSize: number;
	
	    static createFrom(source: any = {}) {
	        return new OrganizePreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourcePath = source["sourcePath"];
	        this.destPath = source["destPath"];
	        this.sourceName = source["sourceName"];
	        this.folderName = source["folderName"];
	        this.fileSize = source["fileSize"];
	    }
	}
	export class OrganizeResult {
	    sourcePath: string;
	    destPath: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new OrganizeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourcePath = source["sourcePath"];
	        this.destPath = source["destPath"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class OrganizeRule {
	    mode: string;
	    target: string;
	    move: boolean;
	    sortInto: string;
	
	    static createFrom(source: any = {}) {
	        return new OrganizeRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.target = source["target"];
	        this.move = source["move"];
	        this.sortInto = source["sortInto"];
	    }
	}
	export class RecycleBinInfo {
	    itemCount: number;
	    size: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new RecycleBinInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.itemCount = source["itemCount"];
	        this.size = source["size"];
	        this.path = source["path"];
	    }
	}
	export class RenamePreview {
	    originalPath: string;
	    newPath: string;
	    originalName: string;
	    newName: string;
	    index: number;
	
	    static createFrom(source: any = {}) {
	        return new RenamePreview(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.originalPath = source["originalPath"];
	        this.newPath = source["newPath"];
	        this.originalName = source["originalName"];
	        this.newName = source["newName"];
	        this.index = source["index"];
	    }
	}
	export class RenameRule {
	    pattern: string;
	    startIndex: number;
	    padding: number;
	    replaceFrom: string;
	    replaceTo: string;
	    fileFilter: string;
	
	    static createFrom(source: any = {}) {
	        return new RenameRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pattern = source["pattern"];
	        this.startIndex = source["startIndex"];
	        this.padding = source["padding"];
	        this.replaceFrom = source["replaceFrom"];
	        this.replaceTo = source["replaceTo"];
	        this.fileFilter = source["fileFilter"];
	    }
	}
	export class ReplaceResult {
	    path: string;
	    matches: number;
	    replaced: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ReplaceResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.matches = source["matches"];
	        this.replaced = source["replaced"];
	        this.error = source["error"];
	    }
	}
	export class SearchResult {
	    path: string;
	    fileName: string;
	    line: number;
	    content: string;
	    fileSize: number;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.fileName = source["fileName"];
	        this.line = source["line"];
	        this.content = source["content"];
	        this.fileSize = source["fileSize"];
	    }
	}

}

export namespace network {
	
	export class BatchPingResult {
	    ip: string;
	    alive: boolean;
	    latency: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchPingResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.alive = source["alive"];
	        this.latency = source["latency"];
	    }
	}
	export class ConnectionInfo {
	    protocol: string;
	    localAddr: string;
	    localPort: number;
	    remoteAddr: string;
	    remotePort: number;
	    state: string;
	    pid: number;
	    process: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.protocol = source["protocol"];
	        this.localAddr = source["localAddr"];
	        this.localPort = source["localPort"];
	        this.remoteAddr = source["remoteAddr"];
	        this.remotePort = source["remotePort"];
	        this.state = source["state"];
	        this.pid = source["pid"];
	        this.process = source["process"];
	    }
	}
	export class DNSResult {
	    hostname: string;
	    type: string;
	    answers: string[];
	    ttl: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new DNSResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.type = source["type"];
	        this.answers = source["answers"];
	        this.ttl = source["ttl"];
	        this.error = source["error"];
	    }
	}
	export class FixResult {
	    action: string;
	    success: boolean;
	    output: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new FixResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.action = source["action"];
	        this.success = source["success"];
	        this.output = source["output"];
	        this.error = source["error"];
	    }
	}
	export class GeoIPResult {
	    ip: string;
	    country: string;
	    region: string;
	    city: string;
	    isp: string;
	    org: string;
	    lat: number;
	    lon: number;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new GeoIPResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.country = source["country"];
	        this.region = source["region"];
	        this.city = source["city"];
	        this.isp = source["isp"];
	        this.org = source["org"];
	        this.lat = source["lat"];
	        this.lon = source["lon"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class LANDevice {
	    ip: string;
	    mac: string;
	    hostname: string;
	    vendor: string;
	    alive: boolean;
	
	    static createFrom(source: any = {}) {
	        return new LANDevice(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.mac = source["mac"];
	        this.hostname = source["hostname"];
	        this.vendor = source["vendor"];
	        this.alive = source["alive"];
	    }
	}
	export class PingResult {
	    target: string;
	    success: boolean;
	    latency: string;
	    ttl: number;
	    sequence: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new PingResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.target = source["target"];
	        this.success = source["success"];
	        this.latency = source["latency"];
	        this.ttl = source["ttl"];
	        this.sequence = source["sequence"];
	        this.error = source["error"];
	    }
	}
	export class PingSummary {
	    target: string;
	    results: PingResult[];
	    sent: number;
	    received: number;
	    lossRate: number;
	    minLatency: string;
	    maxLatency: string;
	    avgLatency: string;
	
	    static createFrom(source: any = {}) {
	        return new PingSummary(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.target = source["target"];
	        this.results = this.convertValues(source["results"], PingResult);
	        this.sent = source["sent"];
	        this.received = source["received"];
	        this.lossRate = source["lossRate"];
	        this.minLatency = source["minLatency"];
	        this.maxLatency = source["maxLatency"];
	        this.avgLatency = source["avgLatency"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PortResult {
	    port: number;
	    state: string;
	    service: string;
	
	    static createFrom(source: any = {}) {
	        return new PortResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.port = source["port"];
	        this.state = source["state"];
	        this.service = source["service"];
	    }
	}
	export class SignalInfo {
	    ssid: string;
	    bssid: string;
	    signal: number;
	    channel: number;
	    auth: string;
	    mhz: number;
	
	    static createFrom(source: any = {}) {
	        return new SignalInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ssid = source["ssid"];
	        this.bssid = source["bssid"];
	        this.signal = source["signal"];
	        this.channel = source["channel"];
	        this.auth = source["auth"];
	        this.mhz = source["mhz"];
	    }
	}
	export class TrafficSample {
	    time: string;
	    download: number;
	    upload: number;
	
	    static createFrom(source: any = {}) {
	        return new TrafficSample(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.download = source["download"];
	        this.upload = source["upload"];
	    }
	}
	export class WiFiProfile {
	    ssid: string;
	    password: string;
	    auth: string;
	
	    static createFrom(source: any = {}) {
	        return new WiFiProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ssid = source["ssid"];
	        this.password = source["password"];
	        this.auth = source["auth"];
	    }
	}

}

export namespace optimize {
	
	export class ChangeResult {
	    name: string;
	    action: string;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ChangeResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.action = source["action"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class CleanResult {
	    target: string;
	    description: string;
	    fileCount: number;
	    freedBytes: number;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CleanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.target = source["target"];
	        this.description = source["description"];
	        this.fileCount = source["fileCount"];
	        this.freedBytes = source["freedBytes"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}
	export class CleanupTarget {
	    path: string;
	    description: string;
	    risk: string;
	    browser?: string;
	    fileCount: number;
	    totalSize: number;
	    checked: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CleanupTarget(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.description = source["description"];
	        this.risk = source["risk"];
	        this.browser = source["browser"];
	        this.fileCount = source["fileCount"];
	        this.totalSize = source["totalSize"];
	        this.checked = source["checked"];
	        this.error = source["error"];
	    }
	}
	export class ContextMenuStatus {
	    installed: boolean;
	    path?: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ContextMenuStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.installed = source["installed"];
	        this.path = source["path"];
	        this.error = source["error"];
	    }
	}
	export class HealthItem {
	    name: string;
	    status: string;
	    value: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new HealthItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.status = source["status"];
	        this.value = source["value"];
	        this.message = source["message"];
	    }
	}
	export class HealthReport {
	    generatedAt: string;
	    score: number;
	    items: HealthItem[];
	    suggestions: string[];
	
	    static createFrom(source: any = {}) {
	        return new HealthReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.generatedAt = source["generatedAt"];
	        this.score = source["score"];
	        this.items = this.convertValues(source["items"], HealthItem);
	        this.suggestions = source["suggestions"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class HostsEntry {
	    ip: string;
	    hostname: string;
	    comment: string;
	    enabled: boolean;
	    line?: string;
	
	    static createFrom(source: any = {}) {
	        return new HostsEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ip = source["ip"];
	        this.hostname = source["hostname"];
	        this.comment = source["comment"];
	        this.enabled = source["enabled"];
	        this.line = source["line"];
	    }
	}
	export class ServiceRecommendation {
	    name: string;
	    action: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceRecommendation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.action = source["action"];
	    }
	}
	export class OptimizationProfile {
	    name: string;
	    description: string;
	    services: ServiceRecommendation[];
	
	    static createFrom(source: any = {}) {
	        return new OptimizationProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.services = this.convertValues(source["services"], ServiceRecommendation);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RegistryScanResult {
	    key: string;
	    value: string;
	    issue: string;
	    category: string;
	    backupPath?: string;
	
	    static createFrom(source: any = {}) {
	        return new RegistryScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	        this.issue = source["issue"];
	        this.category = source["category"];
	        this.backupPath = source["backupPath"];
	    }
	}
	export class RestorePointInfo {
	    name: string;
	    createdAt: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new RestorePointInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.createdAt = source["createdAt"];
	        this.description = source["description"];
	    }
	}
	export class ServiceBackup {
	    name: string;
	    startType: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceBackup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.startType = source["startType"];
	        this.status = source["status"];
	    }
	}
	export class ServiceInfo {
	    name: string;
	    displayName: string;
	    description: string;
	    status: string;
	    startType: string;
	    recommended: string;
	
	    static createFrom(source: any = {}) {
	        return new ServiceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.displayName = source["displayName"];
	        this.description = source["description"];
	        this.status = source["status"];
	        this.startType = source["startType"];
	        this.recommended = source["recommended"];
	    }
	}
	
	export class StartupItem {
	    name: string;
	    command: string;
	    location: string;
	    publisher: string;
	    enabled: boolean;
	    impact: string;
	
	    static createFrom(source: any = {}) {
	        return new StartupItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.command = source["command"];
	        this.location = source["location"];
	        this.publisher = source["publisher"];
	        this.enabled = source["enabled"];
	        this.impact = source["impact"];
	    }
	}
	export class UpdateInfo {
	    name: string;
	    kb: string;
	    status: string;
	    installDate: string;
	    size: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.kb = source["kb"];
	        this.status = source["status"];
	        this.installDate = source["installDate"];
	        this.size = source["size"];
	    }
	}

}

export namespace process {
	
	export class ProcessInfo {
	    pid: number;
	    name: string;
	    cpu: string;
	    memory: string;
	    status: string;
	    user: string;
	    command: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.cpu = source["cpu"];
	        this.memory = source["memory"];
	        this.status = source["status"];
	        this.user = source["user"];
	        this.command = source["command"];
	    }
	}

}

export namespace report {
	
	export class CPUInfo {
	    name: string;
	    cores: number;
	    logicalCores: number;
	    baseClock: string;
	
	    static createFrom(source: any = {}) {
	        return new CPUInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.cores = source["cores"];
	        this.logicalCores = source["logicalCores"];
	        this.baseClock = source["baseClock"];
	    }
	}
	export class DiskInfo {
	    label: string;
	    fileSystem: string;
	    total: number;
	    free: number;
	    usage: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.fileSystem = source["fileSystem"];
	        this.total = source["total"];
	        this.free = source["free"];
	        this.usage = source["usage"];
	    }
	}
	export class MemInfo {
	    total: number;
	    available: number;
	    usage: number;
	
	    static createFrom(source: any = {}) {
	        return new MemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.available = source["available"];
	        this.usage = source["usage"];
	    }
	}
	export class NetworkAdapter {
	    name: string;
	    ip: string;
	    mac: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new NetworkAdapter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.ip = source["ip"];
	        this.mac = source["mac"];
	        this.type = source["type"];
	    }
	}
	export class NetworkInfo {
	    adapters: NetworkAdapter[];
	
	    static createFrom(source: any = {}) {
	        return new NetworkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.adapters = this.convertValues(source["adapters"], NetworkAdapter);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OSInfo {
	    name: string;
	    version: string;
	    buildNumber: string;
	    architecture: string;
	    uptime: string;
	    hostname: string;
	    userName: string;
	
	    static createFrom(source: any = {}) {
	        return new OSInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.buildNumber = source["buildNumber"];
	        this.architecture = source["architecture"];
	        this.uptime = source["uptime"];
	        this.hostname = source["hostname"];
	        this.userName = source["userName"];
	    }
	}
	export class SystemReport {
	    generatedAt: string;
	    os: OSInfo;
	    cpu: CPUInfo;
	    memory: MemInfo;
	    disks: DiskInfo[];
	    network: NetworkInfo;
	    processCount: number;
	
	    static createFrom(source: any = {}) {
	        return new SystemReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.generatedAt = source["generatedAt"];
	        this.os = this.convertValues(source["os"], OSInfo);
	        this.cpu = this.convertValues(source["cpu"], CPUInfo);
	        this.memory = this.convertValues(source["memory"], MemInfo);
	        this.disks = this.convertValues(source["disks"], DiskInfo);
	        this.network = this.convertValues(source["network"], NetworkInfo);
	        this.processCount = source["processCount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace scheduler {
	
	export class TaskInfo {
	    name: string;
	    action: string;
	    time: string;
	    enabled: boolean;
	    taskPath?: string;
	
	    static createFrom(source: any = {}) {
	        return new TaskInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.action = source["action"];
	        this.time = source["time"];
	        this.enabled = source["enabled"];
	        this.taskPath = source["taskPath"];
	    }
	}

}

export namespace screenshot {
	
	export class CaptureResult {
	    path: string;
	    width: number;
	    height: number;
	    size: number;
	    success: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CaptureResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.size = source["size"];
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}

}

export namespace security {
	
	export class EncryptResult {
	    inputPath: string;
	    outputPath: string;
	    success: boolean;
	    error?: string;
	    algorithm?: string;
	    fileSize?: number;
	    encryptedSize?: number;
	
	    static createFrom(source: any = {}) {
	        return new EncryptResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inputPath = source["inputPath"];
	        this.outputPath = source["outputPath"];
	        this.success = source["success"];
	        this.error = source["error"];
	        this.algorithm = source["algorithm"];
	        this.fileSize = source["fileSize"];
	        this.encryptedSize = source["encryptedSize"];
	    }
	}
	export class PasswordResult {
	    password: string;
	    strength: string;
	    length: number;
	
	    static createFrom(source: any = {}) {
	        return new PasswordResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.password = source["password"];
	        this.strength = source["strength"];
	        this.length = source["length"];
	    }
	}
	export class ShredResult {
	    path: string;
	    success: boolean;
	    passes: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ShredResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.success = source["success"];
	        this.passes = source["passes"];
	        this.error = source["error"];
	    }
	}

}

export namespace system {
	
	export class ActivationInfo {
	    activationStatus: string;
	    productId: string;
	    edition: string;
	
	    static createFrom(source: any = {}) {
	        return new ActivationInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.activationStatus = source["activationStatus"];
	        this.productId = source["productId"];
	        this.edition = source["edition"];
	    }
	}
	export class ActivationTool {
	    name: string;
	    description: string;
	    url: string;
	    openSource: boolean;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new ActivationTool(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.url = source["url"];
	        this.openSource = source["openSource"];
	        this.type = source["type"];
	    }
	}
	export class CPUInfo {
	    name: string;
	    cores: number;
	    logicalCores: number;
	    baseClock: string;
	    usage: number;
	    temperature: number;
	
	    static createFrom(source: any = {}) {
	        return new CPUInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.cores = source["cores"];
	        this.logicalCores = source["logicalCores"];
	        this.baseClock = source["baseClock"];
	        this.usage = source["usage"];
	        this.temperature = source["temperature"];
	    }
	}
	export class DiskInfo {
	    label: string;
	    fileSystem: string;
	    total: number;
	    used: number;
	    free: number;
	    usage: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.label = source["label"];
	        this.fileSystem = source["fileSystem"];
	        this.total = source["total"];
	        this.used = source["used"];
	        this.free = source["free"];
	        this.usage = source["usage"];
	    }
	}
	export class HardwareMonitor {
	    cpuUsage: number;
	    memoryUsage: number;
	    memoryUsed: number;
	    memoryTotal: number;
	    diskIO: number;
	    netDown: number;
	    netUp: number;
	    uptime: string;
	
	    static createFrom(source: any = {}) {
	        return new HardwareMonitor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpuUsage = source["cpuUsage"];
	        this.memoryUsage = source["memoryUsage"];
	        this.memoryUsed = source["memoryUsed"];
	        this.memoryTotal = source["memoryTotal"];
	        this.diskIO = source["diskIO"];
	        this.netDown = source["netDown"];
	        this.netUp = source["netUp"];
	        this.uptime = source["uptime"];
	    }
	}
	export class KMSMethod {
	    title: string;
	    steps: string[];
	    command?: string;
	
	    static createFrom(source: any = {}) {
	        return new KMSMethod(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.steps = source["steps"];
	        this.command = source["command"];
	    }
	}
	export class MemoryInfo {
	    total: number;
	    used: number;
	    available: number;
	    usage: number;
	
	    static createFrom(source: any = {}) {
	        return new MemoryInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.used = source["used"];
	        this.available = source["available"];
	        this.usage = source["usage"];
	    }
	}
	export class NetworkAdapter {
	    name: string;
	    type: string;
	    ip: string;
	    mac: string;
	    isConnected: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NetworkAdapter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.ip = source["ip"];
	        this.mac = source["mac"];
	        this.isConnected = source["isConnected"];
	    }
	}
	export class NetworkInfo {
	    adapters: NetworkAdapter[];
	    downloadSpeed: number;
	    uploadSpeed: number;
	
	    static createFrom(source: any = {}) {
	        return new NetworkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.adapters = this.convertValues(source["adapters"], NetworkAdapter);
	        this.downloadSpeed = source["downloadSpeed"];
	        this.uploadSpeed = source["uploadSpeed"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OSInfo {
	    name: string;
	    version: string;
	    buildNumber: string;
	    architecture: string;
	    installDate: string;
	    uptime: string;
	
	    static createFrom(source: any = {}) {
	        return new OSInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.version = source["version"];
	        this.buildNumber = source["buildNumber"];
	        this.architecture = source["architecture"];
	        this.installDate = source["installDate"];
	        this.uptime = source["uptime"];
	    }
	}
	export class PowerPlan {
	    guid: string;
	    name: string;
	    active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PowerPlan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.guid = source["guid"];
	        this.name = source["name"];
	        this.active = source["active"];
	    }
	}
	export class SystemInfo {
	    os: OSInfo;
	    cpu: CPUInfo;
	    memory: MemoryInfo;
	    disks: DiskInfo[];
	    network: NetworkInfo;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = this.convertValues(source["os"], OSInfo);
	        this.cpu = this.convertValues(source["cpu"], CPUInfo);
	        this.memory = this.convertValues(source["memory"], MemoryInfo);
	        this.disks = this.convertValues(source["disks"], DiskInfo);
	        this.network = this.convertValues(source["network"], NetworkInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace upload {
	
	export class UploadResult {
	    success: boolean;
	    fileName: string;
	    serverUrl: string;
	    statusCode: number;
	    response: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new UploadResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.fileName = source["fileName"];
	        this.serverUrl = source["serverUrl"];
	        this.statusCode = source["statusCode"];
	        this.response = source["response"];
	        this.error = source["error"];
	    }
	}

}

