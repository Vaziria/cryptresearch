import torch

a = torch.tensor([2., 3.], requires_grad=True)
b = torch.tensor([6., 4.], requires_grad=True)

Q = 3*a**3 - b**2

print(Q)

external_grad = torch.tensor([1., 1.])
Q.backward(gradient=external_grad)

print(external_grad)
print(a.grad)