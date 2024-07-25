Feature('Stock Management');

Scenario('Verificar la visualización de los insumos existentes', ({ I }) => {
  I.amOnPage('/');
  I.wait(3);

  // Verificar que la tabla de insumos está presente
  I.seeElement('table');
  I.seeElement('table tbody');
  I.seeElement('table tbody tr');
  
  // Verificar la presencia de insumos específicos en la tabla
  I.see('Harina', 'table tbody');
  I.see('10000', 'table tbody');
  I.see('gr', 'table tbody');

  // Asegúrate de que la tabla tenga filas visibles
  I.grabNumberOfVisibleElements('table tbody tr').then((count) => {
    assert(count > 0, 'No se encontraron filas en la tabla.');
  });
});


Scenario('Agregar un nuevo insumo', ({ I }) => {
  I.amOnPage('/');
  I.wait(3);
  I.grabNumberOfVisibleElements('table tbody tr').then((initialCount) => {
    I.fillField('Nombre', 'Manteca');
    I.fillField('Cantidad', 10);
    I.fillField('Unidad', 'kg');
    I.click('Agregar');
    I.waitForFunction((initialCount) => {
      const currentCount = document.querySelectorAll('table tbody tr').length;
      return currentCount === initialCount + 1;
    }, [initialCount], 15);
    I.grabNumberOfVisibleElements('table tbody tr').then((finalCount) => {
      assert.equal(finalCount, initialCount + 1);
    });
  });
});


Scenario('Actualizar un insumo existente', ({ I }) => {
  I.amOnPage('/');
  I.wait(3);
  I.fillField('Nombre', 'Manteca');
  I.fillField('Cantidad', 5);
  I.fillField('Unidad', 'kg');
  I.click('Agregar');
  I.waitForElement('table tbody tr:last-child', 15);
  I.see('Manteca', 'table tbody tr:last-child');
  I.see('15', 'table tbody tr:last-child');
});

Scenario('Descontar un insumo existente', ({ I }) => {
  I.amOnPage('/');
  I.wait(3);
  I.fillField('Nombre', 'Manteca');
  I.fillField('Cantidad', 5);
  I.fillField('Unidad', 'kg');
  I.click('Descontar');
  I.waitForElement('table tbody tr:last-child', 15);
  I.see('Manteca', 'table tbody tr:last-child');
  I.see('10', 'table tbody tr:last-child');
});

Scenario('Eliminar un insumo', ({ I }) => {
  I.amOnPage('/');
  I.wait(3);
  I.grabNumberOfVisibleElements('table tbody tr').then((initialCount) => {
    I.click('table tbody tr:last-child button');
    I.waitForFunction((initialCount) => {
      const currentCount = document.querySelectorAll('table tbody tr').length;
      return currentCount === initialCount - 1;
    }, [initialCount], 15);
    I.grabNumberOfVisibleElements('table tbody tr').then((finalCount) => {
      assert.equal(finalCount, initialCount - 1);
    });
  });
});
