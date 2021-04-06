const app = Vue.createApp({
    data() {
        return {
            activeTab: "home",
            extraSettingsOpen: false,
            platformConfigOpen: [false, ''],
            instances: [
                {
                    "Name": "",
                    "ImgFolders": [
                        ""
                    ],
                    "TimeToQueue": "",
                    "PostInterval": "",
                    "PostDelayAtStartup": "",
                    "Platforms": {
                        "facebook": "",
                        "twitter": ""
                    }
                }
            ]
        }
    },
    mounted() {
        fetch('./config')
            .then(function (response) { return response.json() })
            .then(function (data) { vm.instances = data })
            .then(function () {
                document.querySelector('#home').nextElementSibling.classList.add("after-active")})
            .then(function() {
                vm.instances[0].Picture = "./coolcars.webp"
                vm.instances[1].Picture = "./plantgarden.webp"
                vm.instances[2].Picture = "./johnfunnypic.webp"
            })
        setInterval(function(){
            fetch('./config')
            .then(function (response) { return response.json() })
            .then(function (data) {
                vm.instances = data
             })
            .then(function() {
                vm.instances[0].Picture = "./coolcars.webp"
                vm.instances[1].Picture = "./plantgarden.webp"
                vm.instances[2].Picture = "./johnfunnypic.webp"
            })
        }, 2000)
    },
    computed: {
        tabs() {
            var items = ["home"];
            this.instances.forEach(function (item, i) {
                items.push(item)
            })
            return items
        },
        lastTabActive() {
            return this.instances.length - 1 === this.activeTab
        },
        tabAfterActive() {
            var activeTabIndex = 0
            this.tabs.forEach(function (tab, index) {
                if (vm.activeTab === tab) {
                    activeTabIndex = index
                }
            })
            return this.tabs[activeTabIndex + 1]
        },
        tabBeforeActive() {
            var activeTabIndex = 0
            this.tabs.forEach(function (tab, index) {
                if (vm.activeTab === tab) {
                    activeTabIndex = index
                }
            })
            return this.tabs[activeTabIndex - 1]
        }
    },
    methods: {
        printinfo() {
            console.log(this.instances[this.activeTab])
        },
        nextPostTime(instanceIndex) {
            var d = new Date(this.instances[instanceIndex].NextPostTime)
            return d
        },
        setActive(instanceIndex) {
            this.extraSettingsOpen = false
            this.platformConfigOpen.splice(0, 1, false)

            // nextTick here prevents the exit transition from firing on platformConfig screen when changing
            // the Active instance, because when a new item is set active, new data appears in the platformConfig
            // screen while the exit transition is playing, not very pretty
            this.$nextTick(function () {
                this.activeTab = instanceIndex;

                var items = document.getElementById("sidebar").children;
                for (var i = 0; i < items.length; i++) {
                    items[i].classList.remove("before-active");
                    items[i].classList.remove("active");
                    items[i].classList.remove("after-active");
                }

                if (instanceIndex === 'home') {
                    document.querySelector("#home").previousElementSibling.classList.add("before-active")
                    document.querySelector("#home").nextElementSibling.classList.add("after-active")
                    document.querySelector("#home").classList.add("active")

                } else {
                    document.querySelector("#instance" + instanceIndex).previousElementSibling.classList.add("before-active")
                    document.querySelector("#instance" + instanceIndex).nextElementSibling.classList.add("after-active")
                    document.querySelector("#instance" + instanceIndex).classList.add("active")
                }
            })
        },
        timeUnitFromString(timestring) {
            var numAndUnit = [timestring.slice(0, -1), timestring.slice(-1)]
            console.log(numAndUnit)
            return numAndUnit
        },
        // need to use .splice for all array-changing because Vue has trouble reacting to array changes
        // .splice re-creates the whole array which lets Vue know to re-render the whole array
        // it sounds computationally intensive but is no big deal for small text like folder paths
        newFolder() {
            this.instances[this.activeTab].ImgFolders.splice(this.instances[this.activeTab].ImgFolders.length, 1, "")
        },
        updateFolder(index, event) {
            this.instances[this.activeTab].ImgFolders.splice(index, 1, event.target.value)
        },
        removeFolder(index) {
            this.instances[this.activeTab].ImgFolders.splice(index, 1)
        },
        newInstance() {
            this.instances.push(
            {
                "Name": "Testing",
                "ImgFolders": [
                    ""
                ],
                "TimeToQueue": "2 seconds",
                "PostInterval": "5 hours",
                "PostDelayAtStartup": "random",
                "Platforms": {
                    "facebook": "haha",
                    "twitter": "haha"
                },
                "Status": "needs-configuring",
                "Picture": "johnfunnypic.webp"
            })

            this.$nextTick(function () {
                this.activeTab = this.instances.length - 1;
                var items = document.getElementById("sidebar").children;
                for (var i = 0; i < items.length; i++) {
                    items[i].classList.remove("before-active");
                    items[i].classList.remove("active");
                    items[i].classList.remove("after-active");
                }
    
                document.querySelector(".sidebar-item:nth-last-child(2)").classList.add("active")
                document.querySelector(".active").previousElementSibling.classList.add("before-active")
                document.querySelector(".active").nextElementSibling.classList.add("after-active")

            })

        }
    }
})

app.component("platform-config", {
    props: ['instance', 'platform'],
    template: //html
    `<div id="platform-configure" class="settings-view">
        <div class="heading">
            <div></div>
            <span>{{ titleCase(platform) }} setup for {{ instance }}</span>
            <svg @click="$emit('platform-config-close')" stroke="black" stroke-width="100" style="overflow: visible;" width="20px" viewBox="0 0 305 305" xmlns="http://www.w3.org/2000/svg" fill-rule="evenodd" clip-rule="evenodd" stroke-linejoin="round" stroke-linecap="round" stroke-miterlimit="2">
                <path d="M305 0L0 305M0 0l305 305" fill="none"/>
            </svg>
        </div>

        <div id="platformConfigContent">
            hello
        </div>
    </div>`,
    methods: {
       titleCase(str) {
            return str.slice(0,1).toUpperCase() + str.slice(1)
        }
    }
})

function alertValue() {
    var txt = document.getElementById("test").value;

    fetch('/form', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: "test=" + txt,
    })
}

function getReceive() {
    fetch("/receive").then(function (response) {
        response.text().then(function (txt) {
            vm.receive = txt
        });
    });
}



const vm = app.mount('#app')