# rbxPrivateLimitedInventoryScanner
Shows the collectible/limited item inventory for users with private inventories, which normally is inaccessible.

To start, simply run the program. You will be prompted for either username/userid. Private/unprivate inventory status does not matter, both will work.

Roblox collectible limited items and value are pulled from rolimons.com.

To do list:

-Roblox API changes have rate-limited the "inventory.roblox.com/v1/(userid)/items/asset(itemid)" API, meaning that this program will not work in its current state. In a future build, proxy support will be added. This also has the benefit of allowing terminated-user inventory scanning, which otherwise was unfeasible and unincluded due to the same reasoning.

-Roblox API has also *potentially* deprecated some endpoints used, fixing this should be trivial.

-Error handling is not included as of now, if something has gone wrong, it is most likely due to a change in one of the API endpoints used.

-Ways to export data rather than command line output (probably just write to file if user wants).
