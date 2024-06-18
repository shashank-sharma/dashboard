<script>
    import { Divider } from "attractions";
    import { Icon, CommandLine, Squares2x2 } from "svelte-hero-icons";
    let active = true;
    let activeClass = "sidenav-expand";

    function toggleSideNav() {
        console.log("Toggle");
        if (active) {
            activeClass = "sidenav-collapse";
        } else {
            activeClass = "sidenav-expand";
        }
        active = !active;
    }
</script>

<div class="sidenav {activeClass}">
    <div class="menu-container" on:click={() => toggleSideNav()}>
        <div class="menu-item"><Icon src={Squares2x2} size="32" solid /></div>
    </div>
    <Divider />
    <div class="menu-container">
        <a href="/about" class="menu-item"
            ><Icon src={CommandLine} size="32" /></a
        >
    </div>
    <a href="/settings" class="menu-item">Settings</a>
</div>

<style>
    .sidenav {
        width: 80px;
        background-color: #ffffff;
        color: #fff;
        padding-top: 20px;
        display: flex;
        flex-direction: column;
        transition: width 1.5s;
        box-shadow: 2px 0px 10px rgb(219, 219, 219);
        z-index: 1;
    }

    .sidenav-expand {
        display: block;
    }

    .sidenav-collapse {
        display: none;
    }

    .menu-container {
        display: flex;
        border-right: 3px solid transparent;
        align-items: center;
    }

    .menu-item {
        padding: 10px;
        cursor: pointer;
        color: black;
        display: flex;
        margin-inline: auto;
    }

    .menu-container:hover {
        border-right-color: black;
    }

    .icon {
        margin-right: 10px;
    }
</style>
