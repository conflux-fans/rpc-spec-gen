fs = require('fs');

ethereum = require("../ethereum.json")
function getNameInfos() {
    let mConfig = {}
    ethereum.methods.map(m => mConfig[m.name] = {
        Summary: m.summary,
        ParamNames: getParamNames(m.params),
        ResultName: m.result.name,
        Description: m.result.description
    })

    let result = Object.keys(mConfig).map((k) => {
        v = mConfig[k]
        let res = `"${k}":{`
        if (v.Summary) {
            res += `\n\tSummary: "${mConfig[k].Summary}",`
        }
        if (v.ParamNames.length > 0) {
            const paramNames = "[]string{" + v.ParamNames.map(name => `"${name}"`).join(',') + "}"
            res += `\n\tParamNames: ${paramNames},`
        }
        if (v.ResultName) {
            res += `\n\tResultName: "${mConfig[k].ResultName}",`
        }
        if (v.Description) {
            res += `\n\tDescription: "${mConfig[k].Description}",`
        }
        res += `\n},`
        return res
    }).join('\n')
    return result
}

function getParamNames(params) {
    return params.map(p => {
        if (p.name) {
            return p.name
        }
        if (p["$ref"]) {
            return p["$ref"].split('/').pop().toLowerCase()
        }
        throw new Error("not find name definition for " + JSON.stringify(p))
    })
}

fs.writeFileSync("./out.json", getNameInfos())