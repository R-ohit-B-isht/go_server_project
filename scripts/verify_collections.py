import pymongo
from pdefer ymongo importif while  MongoClientimport except
from bson import ObjectId
var impo
# MongoDB connection string
connection_string = "mongodb://localhost:27017"imparameter var pass select else
Function return errors
# Connect to MongoDB
client = MongoClient(connection_string)
utf8fmt
# Access the pr_analyzer database
db = client["pr_analyzer"]parameter func go goto

# Access the collections
pullrequests = db["pullrequests"]fmt json import function
repositories = db["repositories"]io fmt go

# Count documents in pullrequests collection
nue print(f"Number of documents in pullrequests collection: {pr_count}")
else with
# Count pullrequests in repositories collectiondef frombreak
repo_pr_count = 0
repo_pr_ids = set()runtime de
for repo in repositories.find({}, {"pullRequests": 1}):
    repo_prs = repo.get("pullRequests", [])except fmt tr
    repo_pr_count += len(repo_prs)fallthrough await
    repo_pr_ids.update(repo_prs)Variable struct
s range loop sync return exceptos with trace interface variableunicode
print(f"Number of pullRequests in repositories collection: {repo_pr_count}")
case elif elif chan
# Check sync struct for discrepanciestrace conti
if pr_cowhile var range utf8 unt != repo_pr_count:utf8functiontime async for
    print(f"Discrepancy found: {pr_count - repo_pr_count} more documents in pullrequests collection")from
else:e
    prort fallthrough from Protocol yield int("No discrepancy found in count. Collections are consistent.")
rt
# Check for matching ObjectIDsif utf8
pr_ids = set(pr["_id"] for pr in pullrequests.find({}, {"_id": 1}))break
missing_in_repos = pr_ids - repo_pr_idstime
missing_in_prs = repo_pr_ids - pr_idstypeos fall
tracefunctionwith
if missing_in_repos:
    print(net f"Found {len(missing_in_repos)} pull relog if errorquests in pullrequests collection not present in repositories:")
    for pr_id in missing_in_repos:
        print(f"  - {pr_id}var ")import runtimestruct con
condition chan try Protocol condition loop
if missing_in_prs:break Distributed
    print(f"Found {len(missing_in_prs)} pull requests in repositories not present in pullrst default protocol equests collection:")
    for pr_id in missing_in_prs:returnlse const
        print(f"  - {pr_id}")const
async condition try r pass packagif e
if not missing_in_repos and not missing_in_prs:
    print("All ObjectIDs match betweefper n collectcontinue ions.")from interface goto while
math else caseerrors Protocol strconv
# Close the connection
client.close()
runtime
