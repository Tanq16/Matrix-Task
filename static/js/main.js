// DOM Elements
const modal = document.getElementById('addTaskModal');
const addTaskForm = document.getElementById('addTaskForm');
const taskContent = document.getElementById('taskContent');
const taskQuadrant = document.getElementById('taskQuadrant');

// Utility Functions
const showNotification = (message, isError = false) => {
    const notification = document.createElement('div');
    notification.className = `notification ${isError ? 'error' : 'success'}`;
    notification.textContent = message;
    document.body.appendChild(notification);

    // Animate and remove notification
    setTimeout(() => {
        notification.classList.add('fade-out');
        setTimeout(() => notification.remove(), 300);
    }, 3000);
};

const createTaskElement = (task) => {
    const taskElement = document.createElement('div');
    taskElement.className = 'task-card';
    taskElement.dataset.taskId = task.id;
    
    taskElement.innerHTML = `
        <p class="task-content">${escapeHtml(task.content)}</p>
        <div class="task-actions">
            <button class="complete-btn" onclick="completeTask('${task.id}')">✓</button>
            <button class="delete-btn" onclick="deleteTask('${task.id}')">×</button>
        </div>
    `;
    
    return taskElement;
};

const escapeHtml = (unsafe) => {
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
};

// API Calls
const api = {
    addTask: async (content, quadrant) => {
        try {
            const response = await fetch('/api/tasks', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ content, quadrant }),
            });
            
            if (!response.ok) throw new Error('Failed to add task');
            const result = await response.json();
            console.log('Add task response:', result);
            return result;
        } catch (error) {
            console.error('Error adding task:', error);
            throw error;
        }
    },

    completeTask: async (taskId) => {
        try {
            const response = await fetch('/api/tasks/complete', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ id: taskId }),
            });
            
            if (!response.ok) throw new Error('Failed to complete task');
            const result = await response.json();
            console.log('Complete task response:', result);
            return result;
        } catch (error) {
            console.error('Error completing task:', error);
            throw error;
        }
    },

    deleteTask: async (taskId) => {
        try {
            const response = await fetch(`/api/tasks?id=${taskId}`, {
                method: 'DELETE',
            });
            
            if (!response.ok) throw new Error('Failed to delete task');
            return await response.json();
        } catch (error) {
            console.error('Error deleting task:', error);
            throw error;
        }
    },
};

// Modal Functions
window.showAddTaskForm = (quadrant) => {
    modal.style.display = 'flex';
    taskQuadrant.value = quadrant;
    taskContent.focus();
    
    // Add click outside listener
    modal.addEventListener('click', (e) => {
        if (e.target === modal) hideAddTaskForm();
    });
};

window.hideAddTaskForm = () => {
    modal.style.display = 'none';
    addTaskForm.reset();
};

// Task Operations
window.completeTask = async (taskId) => {
    try {
        const result = await api.completeTask(taskId);
        console.log('Task completed:', result);
        
        // Find and remove the task element with animation
        const taskElement = document.querySelector(`[data-task-id="${taskId}"]`);
        if (taskElement) {
            const quadrant = taskElement.closest('.matrix-quadrant');
            const taskList = quadrant.querySelector('.task-list');
            
            taskElement.classList.add('fade-out');
            
            setTimeout(() => {
                // Remove the task element
                taskElement.remove();
                
                // Check if there are any remaining tasks in this quadrant
                const remainingTasks = taskList.querySelectorAll('.task-card');
                if (remainingTasks.length === 0) {
                    // If no tasks remain, add the empty state
                    const emptyState = document.createElement('div');
                    emptyState.className = 'empty-state';
                    emptyState.innerHTML = '<p>No tasks in this quadrant</p>';
                    taskList.appendChild(emptyState);
                }
                
                showNotification('Task completed successfully');
            }, 300);
        }
    } catch (error) {
        console.error('Failed to complete task:', error);
        showNotification('Failed to complete task', true);
    }
};

window.deleteTask = async (taskId) => {
    if (!confirm('Are you sure you want to delete this task?')) return;
    
    try {
        await api.deleteTask(taskId);
        
        // Find and remove the task element with animation
        const taskElement = document.querySelector(`[data-task-id="${taskId}"]`);
        if (taskElement) {
            taskElement.classList.add('fade-out');
            setTimeout(() => {
                // Remove the task element
                taskElement.remove();
                
                // Check if there are any remaining tasks in this quadrant
                const remainingTasks = taskList.querySelectorAll('.task-card');
                if (remainingTasks.length === 0) {
                    // If no tasks remain, add the empty state
                    const emptyState = document.createElement('div');
                    emptyState.className = 'empty-state';
                    emptyState.innerHTML = '<p>No tasks in this quadrant</p>';
                    taskList.appendChild(emptyState);
                }
                
                showNotification('Task deleted');
            }, 300);
        }
    } catch (error) {
        console.error('Failed to complete task:', error);
        showNotification('Failed to delete task', true);
    }
};

// Form Submission
addTaskForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const content = taskContent.value.trim();
    const quadrant = parseInt(taskQuadrant.value);
    
    if (!content) {
        showNotification('Please enter a task description', true);
        return;
    }
    
    try {
        const response = await api.addTask(content, quadrant);
        
        // Add new task to the appropriate quadrant
        const quadrantElement = document.querySelector(`[data-quadrant="${quadrant}"] .task-list`);
        if (quadrantElement && response.data) {
            // Remove empty state if it exists
            const emptyState = quadrantElement.querySelector('.empty-state');
            if (emptyState) {
                emptyState.remove();
            }

            const taskElement = createTaskElement(response.data);
            quadrantElement.appendChild(taskElement);
        }
        
        hideAddTaskForm();
        showNotification('Task added successfully');
    } catch (error) {
        showNotification('Failed to add task', true);
    }
});

// Add keyboard shortcuts
document.addEventListener('keydown', (e) => {
    // Escape key closes modal
    if (e.key === 'Escape' && modal.style.display === 'flex') {
        hideAddTaskForm();
    }
    
    // Numbers 1-4 while holding Alt opens add task form for respective quadrant
    if (e.altKey && ['1', '2', '3', '4'].includes(e.key)) {
        e.preventDefault();
        showAddTaskForm(parseInt(e.key));
    }
});

// Add drag and drop functionality
document.addEventListener('DOMContentLoaded', () => {
    let draggedTask = null;
    
    // Add drag events to task cards
    document.querySelectorAll('.task-card').forEach(task => {
        task.draggable = true;
        
        task.addEventListener('dragstart', (e) => {
            draggedTask = task;
            e.dataTransfer.effectAllowed = 'move';
            task.classList.add('dragging');
        });
        
        task.addEventListener('dragend', () => {
            task.classList.remove('dragging');
        });
    });
    
    // Add drop zones to quadrants
    document.querySelectorAll('.matrix-quadrant').forEach(quadrant => {
        quadrant.addEventListener('dragover', (e) => {
            e.preventDefault();
            e.dataTransfer.dropEffect = 'move';
            quadrant.classList.add('drag-over');
        });
        
        quadrant.addEventListener('dragleave', () => {
            quadrant.classList.remove('drag-over');
        });
        
        quadrant.addEventListener('drop', async (e) => {
            e.preventDefault();
            quadrant.classList.remove('drag-over');
            
            if (draggedTask) {
                const newQuadrant = parseInt(quadrant.dataset.quadrant);
                const taskId = draggedTask.dataset.taskId;
                
                try {
                    await api.updateTaskQuadrant(taskId, newQuadrant);
                    quadrant.querySelector('.task-list').appendChild(draggedTask);
                    showNotification('Task moved successfully');
                } catch (error) {
                    showNotification('Failed to move task', true);
                }
            }
        });
    });
});

// Add CSS class for notifications
const style = document.createElement('style');
style.textContent = `
    .notification {
        position: fixed;
        bottom: 20px;
        right: 20px;
        padding: 12px 24px;
        border-radius: 4px;
        background-color: var(--color-success);
        color: white;
        z-index: 1000;
        animation: slideIn 0.3s ease;
    }
    
    .notification.error {
        background-color: var(--color-danger);
    }
    
    .notification.fade-out {
        animation: fadeOut 0.3s ease;
    }
    
    .task-card.dragging {
        opacity: 0.5;
    }
    
    .matrix-quadrant.drag-over {
        background-color: var(--color-primary);
        opacity: 0.8;
    }
    
    @keyframes slideIn {
        from { transform: translateX(100%); }
        to { transform: translateX(0); }
    }
    
    @keyframes fadeOut {
        from { opacity: 1; }
        to { opacity: 0; }
    }
`;
document.head.appendChild(style);

// Add empty state styling
const additionalStyles = `
    .empty-state {
        text-align: center;
        padding: var(--spacing-xl);
        color: var(--color-text-light);
        font-style: italic;
    }

    .task-list:empty {
        min-height: 100px;
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: rgba(255, 255, 255, 0.5);
        border-radius: var(--radius-md);
    }

    .task-list:empty::after {
        content: 'No tasks in this quadrant';
        color: var(--color-text-light);
        font-style: italic;
    }
`;

// Add the new styles to the existing style element
document.querySelector('style').textContent += additionalStyles;
