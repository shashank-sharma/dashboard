<!-- src/lib/components/notifications/NotificationIcon.svelte -->
<script lang="ts">
    import {
        Bell,
        Calendar,
        CheckSquare,
        Brain,
        Mail,
        AlertCircle,
    } from "lucide-svelte";

    export let type:
        | "system"
        | "task"
        | "habit"
        | "focus"
        | "calendar"
        | "email";
    export let size: number = 20;

    const icons = {
        system: AlertCircle,
        task: CheckSquare,
        habit: Brain,
        focus: Brain,
        calendar: Calendar,
        email: Mail,
    };

    $: Icon = icons[type] || Bell;
</script>

<Icon {size} />
