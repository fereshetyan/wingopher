export namespace models {
	
	export class AppData {
	    id: string;
	    category: string;
	    choco: string;
	    content: string;
	    description: string;
	    link: string;
	    winget: string;
	    foss: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.category = source["category"];
	        this.choco = source["choco"];
	        this.content = source["content"];
	        this.description = source["description"];
	        this.link = source["link"];
	        this.winget = source["winget"];
	        this.foss = source["foss"];
	    }
	}

}

