<script lang="ts">
    import { foodLogStore } from "../stores/food-log.store";
    import { FOOD_TAGS, DEFAULT_FOOD_LOG_FORM } from "../types";
    import { onMount } from "svelte";
    import { getLocalTimeZone, today } from "@internationalized/date";
    import type { DateValue } from "@internationalized/date";

    // Check if these UI components exist in the project
    import { Button } from "$lib/components/ui/button";
    import { Label } from "$lib/components/ui/label";
    import { Input } from "$lib/components/ui/input";
    import * as Select from "$lib/components/ui/select";
    import { toast } from "svelte-sonner";
    import { X, Upload } from "lucide-svelte";
    import DateTimePicker from "$lib/components/DateTimePicker.svelte";

    // Make initialTag accept string | null | undefined
    export let initialTag: string | null | undefined = null;
    // Add compact mode for modal display
    export let compact = false;
    // Add onSuccess callback for when form submission is successful
    export let onSuccess: () => void = () => {};

    let formData = { ...DEFAULT_FOOD_LOG_FORM };
    let imagePreview: string | null = null;
    let fileInput: HTMLInputElement;
    let isDragging = false;
    let isSubmitting = false;
    let dropZone: HTMLDivElement;

    // DateValue for the DateTimePicker
    let datePickerValue: DateValue = today(getLocalTimeZone());

    // JavaScript Date for our form data
    let dateValue: Date = new Date();

    // Create a string value for the select component
    let selectedTag = initialTag || formData.tag;

    // Set the initial tag if provided
    $: if (initialTag && initialTag !== formData.tag) {
        selectedTag = initialTag;
        formData.tag = initialTag;
    }

    // Initial setup to ensure the date field in formData is set to ISO string
    onMount(() => {
        // Initialize with the current date/time
        formData.date = dateValue.toISOString();

        // Set initial tag if provided
        if (initialTag) {
            selectedTag = initialTag;
            formData.tag = initialTag;
        }
    });

    // Handle date change from the DateTimePicker
    function handleDateChange(date: Date) {
        if (date instanceof Date && !isNaN(date.getTime())) {
            dateValue = date;
            formData.date = dateValue.toISOString();
        } else {
            console.warn("Invalid date received from DateTimePicker:", date);
        }
    }

    // Reset form after submission
    function resetForm() {
        formData = { ...DEFAULT_FOOD_LOG_FORM };
        dateValue = new Date(); // Reset to current date/time
        datePickerValue = today(getLocalTimeZone()); // Reset DatePicker value
        formData.date = dateValue.toISOString();
        selectedTag = formData.tag;
        imagePreview = null;
        if (fileInput) fileInput.value = "";
    }

    // Handle file selection
    function handleFileSelect(event: Event) {
        const input = event.target as HTMLInputElement;
        if (!input.files?.length) return;

        handleFile(input.files[0]);
    }

    // Process the selected file
    function handleFile(file: File) {
        if (!file) return;

        // Validate file type
        if (!file.type.startsWith("image/")) {
            toast.error("Please select an image file");
            return;
        }

        // Validate file size (5MB max)
        if (file.size > 5 * 1024 * 1024) {
            toast.error("Image size must be less than 5MB");
            return;
        }

        // Create preview and set form data
        const reader = new FileReader();
        reader.onload = (e) => {
            imagePreview = e.target?.result as string;
            formData.image = file;
        };
        reader.readAsDataURL(file);
    }

    // Handle drag and drop events
    function handleDragOver(event: DragEvent) {
        event.preventDefault();
        isDragging = true;
    }

    function handleDragLeave() {
        isDragging = false;
    }

    function handleDrop(event: DragEvent) {
        event.preventDefault();
        isDragging = false;

        if (event.dataTransfer?.files?.length) {
            handleFile(event.dataTransfer.files[0]);
        }
    }

    // Handle select change
    function handleTagChange(tag: string) {
        selectedTag = tag;
        formData.tag = tag;
    }

    // Handle form submission
    async function handleSubmit() {
        if (!formData.name) {
            toast.error("Please enter a name for your food item");
            return;
        }

        if (!formData.tag) {
            toast.error("Please select a tag for your food item");
            return;
        }

        if (!formData.image) {
            toast.error("Please upload an image of your food");
            return;
        }

        // Double-check the date is valid
        if (!formData.date || isNaN(new Date(formData.date).getTime())) {
            toast.error("Please select a valid date and time");
            return;
        }

        isSubmitting = true;

        try {
            const result = await foodLogStore.addEntry(formData);

            if (result.success) {
                toast.success("Food item added successfully!");
                resetForm();
                // Call onSuccess callback if provided
                onSuccess();
            } else {
                toast.error("Failed to add food item");
            }
        } catch (error) {
            console.error("Error adding food item:", error);
            toast.error("An unexpected error occurred");
        } finally {
            isSubmitting = false;
        }
    }

    // Handle image removal
    function removeImage() {
        imagePreview = null;
        formData.image = undefined;
        if (fileInput) fileInput.value = "";
    }

    // Initialize drag and drop listeners
    onMount(() => {
        if (dropZone) {
            dropZone.addEventListener("dragover", handleDragOver);
            dropZone.addEventListener("dragleave", handleDragLeave);
            dropZone.addEventListener("drop", handleDrop);
        }

        return () => {
            if (dropZone) {
                dropZone.removeEventListener("dragover", handleDragOver);
                dropZone.removeEventListener("dragleave", handleDragLeave);
                dropZone.removeEventListener("drop", handleDrop);
            }
        };
    });
</script>

<div class="food-log-upload space-y-4">
    {#if !compact}
        <h2 class="text-2xl font-semibold tracking-tight">
            Add New Food Entry
        </h2>
    {/if}

    <form on:submit|preventDefault={handleSubmit} class="space-y-4">
        <div class="space-y-2">
            <Label for="food-name">Food Name</Label>
            <Input
                id="food-name"
                type="text"
                placeholder="What did you eat?"
                bind:value={formData.name}
                required
            />
        </div>

        <!-- Category selection -->
        <div class="space-y-2">
            <Label for="food-tag">Category</Label>
            <div class="select-wrapper">
                <select
                    id="food-tag"
                    class="w-full p-2 border rounded-md"
                    bind:value={selectedTag}
                    on:change={() => (formData.tag = selectedTag)}
                >
                    {#each FOOD_TAGS as tag}
                        <option value={tag.value}>{tag.label}</option>
                    {/each}
                </select>
            </div>
        </div>

        <!-- Date and time selection -->
        <div class="space-y-2">
            <Label for="food-date">When did you eat this?</Label>
            <DateTimePicker
                bind:value={datePickerValue}
                placeholder="Select date and time"
                on:change={(e) => handleDateChange(e.detail)}
            />
            <div class="text-xs text-muted-foreground mt-1">
                {dateValue ? dateValue.toLocaleString() : "No date selected"}
            </div>
        </div>

        <div class="space-y-2">
            <Label for="food-image">Food Image</Label>
            <div
                bind:this={dropZone}
                class="border-2 border-dashed p-6 rounded-lg text-center cursor-pointer transition-colors duration-200 flex flex-col items-center justify-center"
                class:border-primary={isDragging}
                class:bg-primary-50={isDragging}
                class:border-gray-300={!isDragging}
                on:click={() => fileInput?.click()}
                on:keydown={(e) => e.key === "Enter" && fileInput?.click()}
                tabindex="0"
                role="button"
                aria-label="Upload image"
            >
                {#if imagePreview}
                    <div
                        class="relative max-w-full max-h-64 overflow-hidden mb-2 group"
                    >
                        <img
                            src={imagePreview}
                            alt="Preview"
                            class="max-h-64 object-contain rounded-md"
                        />
                        <button
                            type="button"
                            class="absolute top-2 right-2 bg-white/80 p-1 rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
                            on:click|stopPropagation={removeImage}
                            aria-label="Remove image"
                        >
                            <X class="h-4 w-4" />
                        </button>
                    </div>
                    <span class="text-sm text-gray-500"
                        >Click to change image</span
                    >
                {:else}
                    <Upload class="h-12 w-12 text-gray-400 mb-2" />
                    <p class="text-sm text-gray-500">
                        Drag and drop an image, or click to select
                    </p>
                    <p class="text-xs text-gray-400 mt-1">
                        JPG, PNG, GIF up to 5MB
                    </p>
                {/if}
                <input
                    bind:this={fileInput}
                    type="file"
                    id="food-image"
                    accept="image/*"
                    class="hidden"
                    on:change={handleFileSelect}
                />
            </div>
        </div>

        <Button type="submit" class="w-full" disabled={isSubmitting}>
            {#if isSubmitting}
                Adding...
            {:else}
                Add to Food Log
            {/if}
        </Button>
    </form>
</div>

<style>
    .select-wrapper {
        position: relative;
    }

    .select-wrapper select {
        appearance: none;
        background-color: var(--background);
        color: var(--foreground);
        border-color: var(--input);
        height: 40px;
        padding-left: 14px;
        padding-right: 28px;
    }

    .select-wrapper::after {
        content: "";
        position: absolute;
        right: 14px;
        top: 50%;
        transform: translateY(-50%);
        width: 0;
        height: 0;
        border-left: 5px solid transparent;
        border-right: 5px solid transparent;
        border-top: 5px solid currentColor;
        pointer-events: none;
    }
</style>
