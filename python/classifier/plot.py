import csv
import matplotlib.pyplot as plt
from sklearn.metrics import confusion_matrix, ConfusionMatrixDisplay


labels = []
with open('labels.txt') as fp:
    labels = [line.strip() for line in fp.readlines()]

y_true = []
y_pred = []
with open('results.tsv', 'r') as fp:
    reader = csv.reader(fp, delimiter='\t')
    for row in reader:
        y_true.append(labels[int(row[1])])
        y_pred.append(labels[int(row[2])])

cm = confusion_matrix(y_true, y_pred)
print(cm)

disp = ConfusionMatrixDisplay(confusion_matrix=cm,
                              display_labels=labels)
disp = disp.plot(include_values=True,
                 values_format='3d')
plt.savefig('confusuion_matrix.png')
plt.show()
