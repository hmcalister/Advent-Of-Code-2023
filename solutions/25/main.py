import networkx as nx
import matplotlib.pyplot as plt

with open("puzzleInput", "r") as f:
    lines = f.readlines()

G = nx.Graph()
allNodes = set()

for line in lines:
    currentComponent, line = line.split(":")
    neighborComponents = line.split()
    print(currentComponent, neighborComponents)

    G.add_node(currentComponent)
    allNodes.add(currentComponent)
    for neighbor in neighborComponents:
        G.add_node(neighbor)
        allNodes.add(neighbor)
        G.add_edge(currentComponent, neighbor, capacity=1)

# nx.draw(G)
# plt.show()

allNodes = list(allNodes)
cut, partition = nx.minimum_cut(G, allNodes[0], allNodes[3])
print(f"CUT SIZE: {cut}")
# print(partition)
print(len(partition[0]) * len(partition[1]))