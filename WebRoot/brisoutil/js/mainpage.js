import { createElement } from "/node_modules/@cargowebserver/cargowebcomponents/components/element.js";
import { randomUUID } from "/node_modules/@cargowebserver/cargowebcomponents/components/utility.js";
import { attachAutoComplete } from "/node_modules/@cargowebserver/cargowebcomponents/components/autocomplete/autocomplete.js";
import "/node_modules/@cargowebserver/cargowebcomponents/components/growl/growl.js";
import '/node_modules/@polymer/paper-icon-button/paper-icon-button.js';
import '/node_modules/@polymer/iron-icons/iron-icons.js';

export class MainPage {
    constructor() {

        // Keep track of the current user.
        this.currentUser = null

        // Create the main page.

        // The main panel.
        this.panel = createElement(null, { "tag": "div", "class": "mainpage-panel" })

        // Append the panel in the body
        document.body.appendChild(this.panel.element)

        // Now I will create the layout.
        this.headerPanel = this.panel.appendElement({ "tag": "div", "class": "mainpage-header" }).down()
        this.headerPanel.appendElement({ "tag": "div", "class": "mainpage-title", "innerHtml": "Rapports de Bris d'outils et situation problématique" })

        // The content...
        this.contentPanel = this.panel.appendElement({ "tag": "div", "class": "mainpage-content" }).down()

        this.contentPanel.element.onscroll = function (headerPanel) {
            return function () {
                if (this.scrollTop > 10) {
                    // set the
                    headerPanel.element.style.boxShadow= "0 2px 2px 0 rgba(0, 0, 0, 0.14), 0 1px 5px 0 rgba(0, 0, 0, 0.12), 0 3px 1px -2px rgba(0, 0, 0, 0.2)"
                }else{
                    headerPanel.element.style.boxShadow= ""
                }
            }}(this.headerPanel)

        // Now will append the sections...
        this.form = this.contentPanel.appendElement({ "tag": "div", "class": "form" }).down()

        // the number.
        this.numberInput = this.form.appendElement({ "tag": "div", "class": "form-section" }).down()
            .appendElement({ "tag": "div", "style": "display: table; width: 100%;" }).down()
            .appendElement({ "tag": "div", "style": "display: table-row; width: 100%;" }).down()
            .appendElement({ "tag": "div", "style": "display: table-cell", "innerHtml": "Numéros" })
            .appendElement({ "tag": "div", "style": "display: table-cell; width: 100%;" }).down()
            .appendElement({ "tag": "input", "type": "number" }).down()

        this.numberInput.element.onkeyup = function (mainPage) {
            return function (evt) {
                if (evt.keyCode == 13) {
                    mainPage.setBris()
                }
                if(this.value.length == 0){
                    mainPage.clear()
                }
            }
        }(this)

        // The general info section
        this.generalInfo = this.form.appendElement({ "tag": "div", "class": "form-section" }).down()

        // Append a yes no question.
        function appendYesNoQuestion(section, quesiton, yesText, noText) {
            var yesUuid = randomUUID()
            var noUuid = randomUUID()
            var yesTextUuid = randomUUID()
            var noTextUuid = randomUUID()

            section.appendElement({ "tag": "div", "style": "display: flex; width:100%; flex-direction: column;" }).down()
                .appendElement({ "tag": "div", "class": "form-section-title", "innerHtml": quesiton })
                .appendElement({ "tag": "div", "style": "display: table; width:100%;" }).down()
                .appendElement({ "tag": "div", "style": "display: table-row; width:100%;" }).down()
                .appendElement({"tag":"div", "style":"display: table-cell;"}).down()
                .appendElement({ "tag": "input", "id": yesUuid, "type": "checkbox", "style": "vertical-align: middle; text-align: center;" }).up()
                .appendElement({ "tag": "div","id":yesTextUuid, "style": "display: table-cell;", "innerHtml": yesText }).up()
                .appendElement({ "tag": "div", "style": "display: table-row;" }).down()
                .appendElement({"tag":"div", "style":"display: table-cell;"}).down()
                .appendElement({ "tag": "input", "id": noUuid, "type": "checkbox", "style": "display:vertical-align: middle; ; text-align: center;" }).up()
                .appendElement({ "tag": "div","id":noTextUuid, "style": "display: table-cell;", "innerHtml": noText })

            return { "yesBtn": section.getChildById(yesUuid), "noBtn": section.getChildById(noUuid), "yesText": section.getChildById(yesTextUuid), "noText": section.getChildById(noTextUuid)}
        }

        // Question: 1
        this.isBrokenBtns = appendYesNoQuestion(this.generalInfo, "L'outil est-il brisé?", "oui", "non (outil problématique)")
        this.isBrokenBtns.yesText.element.style.width = this.isBrokenBtns.noText.element.style.width = "100%"

        // Question: 2
        this.hasNcBtns = appendYesNoQuestion(this.generalInfo, "Le bris a-t-il causé une non conformité?", "oui", "non")

        this.hasNcDiv = this.hasNcBtns.yesBtn.parentElement.parentElement
        .appendElement({"tag":"div", "style":"display: table-cell; visibility: hidden; width: 100%;"}).down()

        this.ncTextBox = this.hasNcDiv.appendElement({"tag":"div", "style":"display: table-cell;", "innerHtml":"Nc"})
        .appendElement({"tag":"div", "style":"display: table-cell;"}).down()
        .appendElement({"tag":"input", "style":"display: table-cell;"}).down()
        
        // here I will append the nc input...

        // Question: 3
        this.isArtisBtns = appendYesNoQuestion(this.generalInfo, "Artist s'est-il déclanché?", "oui", "non")
        this.isArtisBtns.yesText.element.style.width = this.isArtisBtns.noText.element.style.width = "100%"

        // The specific informations.
        this.specificInfo = this.form.appendElement({ "tag": "div", "class": "form-section" }).down()
            .appendElement({ "tag": "div", "style": "display: table; width: 100%; border-spacing: 7px;" }).down()

        function appendInputBox(section, text, list) {

            var input = section.appendElement({ "tag": "div", "style": "display: table-row; width: 100%;" }).down()
                .appendElement({ "tag": "div", "style": "display: table-cell", "innerHtml": text })
                .appendElement({ "tag": "div", "style": "display: table-cell;" }).down()
                .appendElement({ "tag": "input", "style": "display: table-cell; min-width: 150px;" }).down()

            if (list != undefined) {
                attachAutoComplete(input, list, null)
            }

            return input
        }

        this.operateurInput = appendInputBox(this.specificInfo, "Opérateur:(prénom, nom)", employeNames.sort())

        this.machineInput = appendInputBox(this.specificInfo, "Machine", machines.sort())

        this.pieceInput = appendInputBox(this.specificInfo, "No. de pièce", serialNumbers.sort())

        this.operationInput = appendInputBox(this.specificInfo, "No. opération")

        this.toolInput = appendInputBox(this.specificInfo, "No. outil", toolNumbers.sort())

        this.programInput = appendInputBox(this.specificInfo, "Programme", programNumbers.sort())

        // Now the decription.
        this.decriptionSection = this.form.appendElement({ "tag": "div", "class": "form-section" }).down()

        // The section text.
        this.decription = this.decriptionSection.appendElement({ "tag": "div", "style": "display: flex; flex-direction: column;" }).down()
            .appendElement({ "tag": "div", "style": "", "innerHtml": "Description", "title": "Informations pouvant aider a comprendre la cause du bris" })
            .appendElement({ "tag": "textarea", "style": "min-height: 100px;" }).down()

        // Tha action section
        this.actionSection = this.form.appendElement({ "tag": "div", "class": "form-section" }).down()

        this.appendActionBtn = this.actionSection.appendElement({ "tag": "div", "style": "display: table; width: 100%;" }).down()
            .appendElement({ "tag": "div", "style": "display: table-row; width: 100%;" }).down()
            .appendElement({ "tag": "div", "style": "display: table-cell", "innerHtml": "Actions" })
            .appendElement({ "tag": "div", "style": "display: table-cell; width: 100%; text-align: right;" }).down()
            .appendElement({ "tag": "paper-icon-button", "icon": "add-circle", "style": "color: #616161;" }).down()

        // contain the list of action.
        this.actionSection = null;

        // Now I will append the onclick actions.
        this.appendActionBtn.element.onclick = function (mainPage) {
            return function () {
                // The action will contain the date and the name of the person who create it.
                mainPage.createAction()
            }
        }(this)

        // Display the welcome message.
        this.displayMessage("Bienvenu dans le formulaire de bris d'outil", 3000)

    }

    /**
     * Display user message.
     * @param {*} message 
     */
    displayMessage(message, delay) {
        var div = document.createElement("div")
        div.style.padding = "10px"
        if (delay != undefined) {
            div.innerHTML = '<growl-element w="350" h="100" shadow="5" delay="' + delay + '" style="display:none;"></growl-element>'
        } else {
            div.innerHTML = '<growl-element w="350" h="100" shadow="5" style="display:none;"></growl-element>'
        }
        div.firstChild.style.backgroundColor = "white"
        div.firstChild.innerHTML = message
        document.body.appendChild(div)
    }

    // Create a new action...
    createAction() {
        if (this.currentUser == null) {
            this.displayMessage("<div style='padding: 10px';>Vous devez vous indentifier pour ajouter une actions!</div>")
        }
    }

    // Set the interface for existing value.
    setBris() {
        var query = "SELECT [product_id],[machine_id],[operation_number] "
        query += ",[nc_number],[serial_number],[program_number],[tool_id],[employe_id] "
        query += ",[tool_is_broke] ,[description],[state] ,[artis_fire] "
        query += "FROM [BrisOutil].[dbo].[Bris] "
        query += "WHERE id=?"

        var q = new Sql.Query()
        q.setQuery(query)
        q.setConnectionid("bris_outil")
        q.setParameters(JSON.stringify([this.numberInput.element.value]))

        var rqst = new Sql.QueryContextRqst()
        rqst.setQuery(q)

        var metadata = { 'custom-header-1': 'value1' };
        var stream = globular.sqlService.queryContext(rqst, metadata);

        // Get the stream and set event on it...
        stream.on('data', function (response) {
            if (response.hasHeader()) {
                var header = response.getHeader()
            } else if (response.hasRows()) {
                results = JSON.parse(response.getRows())

                // Set the interface with the values.
                mainPage.pieceInput.element.title = results[0][0]
                mainPage.machineInput.element.value = results[0][1]
                mainPage.operationInput.element.value = results[0][2]
                mainPage.pieceInput.element.value = results[0][4]
                mainPage.programInput.element.value = results[0][5]

                mainPage.decription.element.value = results[0][9]
                mainPage.toolInput.element.value = results[0][6]
                var operator = employeById[results[0][7].toUpperCase()]
                mainPage.operateurInput.element.value = operator.firstName + " " + operator.lastName

                if(results[0][8] == 1){
                    mainPage.isBrokenBtns.yesBtn.element.checked = true
                    mainPage.isBrokenBtns.noBtn.element.checked = false
                }else{
                    mainPage.isBrokenBtns.yesBtn.element.checked = false
                    mainPage.isBrokenBtns.noBtn.element.checked = true
                }

                if(results[0][3] != null){
                    if(results[0][3].length > 0){
                        mainPage.hasNcBtns.yesBtn.element.checked = true
                        mainPage.hasNcBtns.noBtn.element.checked = false
                        mainPage.ncTextBox.element.style.visibility = "visible"
                        mainPage.ncTextBox.element.value = results[0][3]
                    }else{
                        mainPage.hasNcBtns.yesBtn.element.checked = false
                        mainPage.hasNcBtns.noBtn.element.checked = true
                        mainPage.ncTextBox.element.style.visibility = "hidden"
                        mainPage.ncTextBox.element.value = ""
                    }

                }else{
                    mainPage.hasNcBtns.yesBtn.element.checked = false
                    mainPage.hasNcBtns.noBtn.element.checked = true
                    mainPage.ncTextBox.element.style.visibility = "hidden"
                    mainPage.ncTextBox.element.value = ""
                }

                if(results[0][11] == 1){
                    mainPage.isArtisBtns.yesBtn.element.checked = true
                    mainPage.isArtisBtns.noBtn.element.checked = false
                }else{
                    mainPage.isArtisBtns.yesBtn.element.checked = false
                    mainPage.isArtisBtns.noBtn.element.checked = true
                }
            }
        });

        // Get the results here if the statut is ok.
        stream.on('status', function (status) {
            if (status.code == 0) {

            }
        });

        stream.on('end', function (end) {
            // stream end signal
        });

    }

    // Set the actions.
    setAction() {

    }

    // Reset all values.
    clear() {
        
        this.pieceInput.element.title = ""
        this.machineInput.element.value = ""
        this.operationInput.element.value = ""
        this.pieceInput.element.value = ""
        this.programInput.element.value = ""
        this.decription.element.value = ""
        this.toolInput.element.value = ""
        this.operateurInput.element.value = ""
        this.isBrokenBtns.yesBtn.element.checked = false
        this.isBrokenBtns.noBtn.element.checked = false
        this.hasNcBtns.yesBtn.element.checked = false
        this.hasNcBtns.noBtn.element.checked = false
        this.ncTextBox.element.style.visibility = "hidden"
        this.ncTextBox.element.value = ""
        this.isArtisBtns.yesBtn.element.checked = false
        this.isArtisBtns.noBtn.element.checked = false
    }

}

// set in the global namespace.
window.MainPage = MainPage