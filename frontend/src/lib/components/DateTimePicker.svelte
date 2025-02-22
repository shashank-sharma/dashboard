<!-- src/lib/components/DateTimePicker.svelte -->
<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Calendar } from "$lib/components/ui/calendar";
    import * as Popover from "$lib/components/ui/popover";
    import { ScrollArea } from "$lib/components/ui/scroll-area";
    import { cn } from "$lib/utils";
    import { CalendarDays } from "lucide-svelte";
    import { createEventDispatcher } from "svelte";

    export let value: Date | undefined = undefined;
    export let placeholder = "Pick date and time";

    const dispatch = createEventDispatcher<{
        change: Date;
    }>();

    let isOpen = false;

    // Create arrays for hours and minutes
    const hours = Array.from({ length: 24 }, (_, i) => i);
    const minutes = Array.from({ length: 12 }, (_, i) => i * 5);

    function formatDate(date: Date | undefined): string {
        if (!date) return "";
        return date.toLocaleString("en-US", {
            month: "2-digit",
            day: "2-digit",
            year: "numeric",
            hour: "2-digit",
            minute: "2-digit",
            hour12: false,
        });
    }

    function handleDateSelect(date: Date | undefined) {
        if (!date) return;

        // Preserve existing time if we have it
        if (value) {
            date.setHours(value.getHours(), value.getMinutes());
        }

        value = date;
        dispatch("change", value);
    }

    function handleTimeChange(type: "hour" | "minute", val: number) {
        if (!value) {
            value = new Date();
        }

        const newDate = new Date(value);

        if (type === "hour") {
            newDate.setHours(val);
        } else {
            newDate.setMinutes(val);
        }

        value = newDate;
        dispatch("change", value);
    }
</script>

<Popover.Root bind:open={isOpen}>
    <Popover.Trigger asChild let:builder>
        <Button
            builders={[builder]}
            variant="outline"
            class={cn(
                "w-full justify-start text-left font-normal",
                !value && "text-muted-foreground",
            )}
        >
            <CalendarDays class="mr-2 h-4 w-4" />
            {value ? formatDate(value) : placeholder}
        </Button>
    </Popover.Trigger>
    <Popover.Content class="w-auto p-0">
        <div class="sm:flex">
            <Calendar
                mode="single"
                selected={value}
                onSelect={handleDateSelect}
                initialFocus
            />
            <div
                class="flex flex-col sm:flex-row sm:h-[300px] divide-y sm:divide-y-0 sm:divide-x border-l"
            >
                <ScrollArea class="w-64 sm:w-auto">
                    <div class="flex sm:flex-col p-2 gap-1">
                        {#each hours as hour}
                            <Button
                                size="icon"
                                variant={value && value.getHours() === hour
                                    ? "default"
                                    : "ghost"}
                                class="sm:w-full shrink-0 aspect-square"
                                on:click={() => handleTimeChange("hour", hour)}
                            >
                                {hour.toString().padStart(2, "0")}
                            </Button>
                        {/each}
                    </div>
                </ScrollArea>
                <ScrollArea class="w-64 sm:w-auto">
                    <div class="flex sm:flex-col p-2 gap-1">
                        {#each minutes as minute}
                            <Button
                                size="icon"
                                variant={value && value.getMinutes() === minute
                                    ? "default"
                                    : "ghost"}
                                class="sm:w-full shrink-0 aspect-square"
                                on:click={() =>
                                    handleTimeChange("minute", minute)}
                            >
                                {minute.toString().padStart(2, "0")}
                            </Button>
                        {/each}
                    </div>
                </ScrollArea>
            </div>
        </div>
    </Popover.Content>
</Popover.Root>
