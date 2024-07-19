Feature('Stock Management');

Scenario('Agregar un nuevo insumo', async ({ I }) => {
    I.amOnPage('/');  // Asegúrate de estar en la página correcta

    I.wait(3);
    // Obtener la cantidad inicial de insumos
    const initialItemCount = await I.grabNumberOfVisibleElements('table tbody tr');
  
    // Ejecutar la lógica para agregar un nuevo insumo
    I.fillField('Nombre', 'Manteca');
    I.fillField('Cantidad', 10);
    I.fillField('Unidad', 'kg');
    I.click('Agregar');
  
    // Esperar a que el nuevo insumo se agregue dinámicamente
    await I.waitForFunction(
      (initialCount) => {
        const currentCount = document.querySelectorAll('table tbody tr').length;
        return currentCount === initialCount + 1;
      },
      [initialItemCount], 
      10
    );
  
    // Verificar que el insumo se haya agregado correctamente
    const newItemCount = await I.grabNumberOfVisibleElements('table tbody tr');
    if (newItemCount !== initialItemCount + 1) {
      throw new Error('No se agregó el insumo correctamente');
  }
  
  });

  Scenario('Actualizar un insumo existente', async ({ I }) => {
    I.amOnPage('/');
    I.wait(3);
  
    // Agregar el mismo insumo para actualizar la cantidad
    I.fillField('Nombre', 'Manteca');
    I.fillField('Cantidad', 5);
    I.fillField('Unidad', 'kg');
    I.click('Agregar');
  
    // Verificar que la cantidad ha sido actualizada
    I.waitForElement('table tbody tr:last-child', 15);
    I.see('Manteca', 'table tbody tr:last-child');
    I.see('15', 'table tbody tr:last-child'); // 10 + 5
    I.see('kg', 'table tbody tr:last-child');
  });
  
  Scenario('Descontar un insumo existente', async ({ I }) => {
    I.amOnPage('/');
    I.wait(3);
  
    // Seleccionar el insumo para descontar
    I.fillField('Nombre', 'Manteca');
    I.fillField('Cantidad', 5);
    I.fillField('Unidad', 'kg');
    I.click('Descontar');
  
    // Verificar que la cantidad ha sido descontada
    I.waitForElement('table tbody tr:last-child', 15);
    I.see('Manteca', 'table tbody tr:last-child');
    I.see('10', 'table tbody tr:last-child'); // 15 - 5
    I.see('kg', 'table tbody tr:last-child');
  });
  
  
  Scenario('Eliminar un insumo', async ({ I }) => {
    I.amOnPage('/');
    I.wait(3);

    const initialItemCount = await I.grabNumberOfVisibleElements('table tbody tr');

    // Eliminar el último insumo
    I.click('table tbody tr:last-child button');
  
    await I.waitForFunction(
        (initialCount) => {
          const currentCount = document.querySelectorAll('table tbody tr').length;
          return currentCount === initialCount - 1;
        },
        [initialItemCount], 
        10
      );
    
      // Verificar que el insumo se haya agregado correctamente
      const newItemCount = await I.grabNumberOfVisibleElements('table tbody tr');
      if (newItemCount !== initialItemCount - 1) {
        throw new Error('No se elimino el insumo correctamente');
      }
    });
  