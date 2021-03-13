<form class="p-4 w-full flex flex-col space-y-4">
    <div class="flex flex-col w-full space-y-1">
        <label for="email-input">E-Mail</label>
        <input type="email" class="border-2 border-primary rounded-md p-2" name="email-input" placeholder="john.doe@example.org"/>
    </div>

    <div class="flex flex-col w-full space-y-1">
        <label for="password-input">Password</label>
        <input type="password" class="border-2 border-primary rounded-md p-2" name="password-input" placeholder="********" />
    </div>

    <div class="flex justify-between py-2">
        <div></div>
        <button type="submit" class="py-2 px-4 text-white bg-primary rounded-md hover:bg-primary-dark">Log In</button>
    </div>
</form>
