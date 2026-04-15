export namespace main {
	
	export class Config {
	    gta_path: string;
	    redux_folders: string[];
	    active_redux: string;
	    auto_inject: boolean;
	    run_on_startup: boolean;
	    start_window_box: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gta_path = source["gta_path"];
	        this.redux_folders = source["redux_folders"];
	        this.active_redux = source["active_redux"];
	        this.auto_inject = source["auto_inject"];
	        this.run_on_startup = source["run_on_startup"];
	        this.start_window_box = source["start_window_box"];
	    }
	}

}

