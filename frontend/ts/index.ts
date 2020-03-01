// Muli font
import "typeface-muli";

// require-confirm forms
(function () {
    for (const elem of document.getElementsByClassName("require-confirm")) {
        (elem as HTMLFormElement).addEventListener("submit", ev => {
            if (!confirm("Are you sure you want to delete this item?"))
                ev.preventDefault();
        })
    }
})();

// load the list of tests
(function () {
    const testTables = Array.from(document.getElementsByClassName("tests-list"))
    const toggles = Array.from(document.getElementsByClassName("toggle-tests"))
    const groups = testTables
        .reduce((m, elem) => {
            const e = elem as HTMLDivElement;
            const group = e.getAttribute("data-test-group");
            const toggle = toggles.find(t => t.getAttribute("data-test-group") == group);
            m.set(group, [e, toggle]);
            return m;
        }, new Map<string, [HTMLDivElement, Element]>());
    let opening = 0;
    const doToggle = (table: HTMLDivElement, toggle: Element, force?: boolean) => {
        const current = table.style.maxHeight === ""; // table showing?
        if (force === undefined) {
            force = !current;
        }
        if (force === current) return;
        if (force) {
            toggle.innerHTML = "[-]"
            table.style.maxHeight = "";
            ++opening;
        } else {
            toggle.innerHTML = "[+]"
            table.style.maxHeight = "0";
            --opening;
        }
        if (opening > 0) {
            allToggle.innerHTML = "[-]";
        } else {
            allToggle.innerHTML = "[+]";
        }
    }
    const items = Array.from(groups.values());
    for (const [table, toggle] of items) {
        toggle.addEventListener("click", () => doToggle(table, toggle))
    }

    const allToggle: Element = document.getElementById("toggle-all-tests");
    allToggle.addEventListener("click", ev => {
        const switchOn = allToggle.innerHTML === "[+]";
        for (const [table, toggle] of items) {
            doToggle(table, toggle, switchOn)
        }
    });
})();
